package ws

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// при уходе клиента со страницы соединение закрывается
func acceptHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, *bufio.ReadWriter, error) {

	// проверяем заголовки
	if r.Header.Get("Upgrade") != "websocket" {
		return nil, nil, errors.New("no websocket")
	}
	if r.Header.Get("Connection") != "Upgrade" && r.Header.Get("Connection") != "upgrade" {
		return nil, nil, errors.New("no upgrade")
	}
	k := r.Header.Get("Sec-Websocket-Key")
	if k == "" {
		return nil, nil, errors.New("no Sec-Websocket-Key")
	}

	// вычисляем ответ
	sum := k + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(sum))
	str := base64.StdEncoding.EncodeToString(hash[:])

	// Берем под контроль соединение https://pkg.go.dev/net/http#Hijacker
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("no Hijacker")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, nil, errors.New("no Hijack")
	}

	// формируем ответ
	bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	bufrw.WriteString("Upgrade: websocket\r\n")
	bufrw.WriteString("Connection: Upgrade\r\n")
	bufrw.WriteString("Sec-Websocket-Accept: " + str + "\r\n\r\n")
	bufrw.Flush()

	return conn, bufrw, nil
}

type frame struct {
	isFin   bool
	opCode  byte
	length  uint64
	payload []byte
}

func readFrame(bufrw *bufio.ReadWriter) (frame, error) {

	// сообщение состоит из одного или нескольких фреймов
	var ret_frame frame
	var message []byte
	for {
		// заголовок состоит из 2 — 14 байт
		buf := make([]byte, 2, 12)
		// читаем первые 2 байта
		_, err := bufrw.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return frame{}, errors.New("no read 1")
		}
		if err == nil {
			finBit := buf[0] >> 7  // фрагментированное ли сообщение
			opCode := buf[0] & 0xf // опкод

			ret_frame.isFin = finBit == 1
			ret_frame.opCode = opCode

			maskBit := buf[1] >> 7 // замаскированы ли данные

			// оставшийся размер заголовка
			extra := 0
			if maskBit == 1 {
				extra += 4 // +4 байта маскировочный ключ
			}

			size := uint64(buf[1] & 0x7f)
			if size == 126 {
				extra += 2 // +2 байта размер данных
			} else if size == 127 {
				extra += 8 // +8 байт размер данных
			}

			if extra > 0 {
				// читаем остаток заголовка extra <= 12
				buf = buf[:extra]
				_, err = bufrw.Read(buf)
				if err != nil {
					return frame{}, errors.New("no read 2")
				}

				if size == 126 {
					size = uint64(binary.BigEndian.Uint16(buf[:2]))
					buf = buf[2:] // подвинем начало буфера на 2 байта
				} else if size == 127 {
					size = uint64(binary.BigEndian.Uint64(buf[:8]))
					buf = buf[8:] // подвинем начало буфера на 8 байт
				}
			}

			// маскировочный ключ
			var mask []byte
			if maskBit == 1 {
				// остаток заголовка, последние 4 байта
				mask = buf
			}

			// данные фрейма
			payload := make([]byte, int(size))
			// читаем полностью и ровно size байт
			_, err = io.ReadFull(bufrw, payload)
			if err != nil {
				return frame{}, errors.New("no read 3")
			}

			// размаскировываем данные с помощью XOR
			if maskBit == 1 {
				for i := 0; i < len(payload); i++ {
					payload[i] ^= mask[i%4]
				}
			}

			// складываем фрагменты сообщения
			message = append(message, payload...)

			ret_frame.payload = message
			ret_frame.length = uint64(len(message))

			if opCode == 8 { // фрейм закрытия
				fmt.Println("opCode = 8")
				return ret_frame, errors.New("closed connect")
			} else if finBit == 1 { // конец сообщения
				fmt.Println(string(ret_frame.payload))
				return ret_frame, nil
			}
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}

}

// функция отправки клиенту заглушки с инициализацией WS
func Ws_handler_start(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	fmt.Println(*r)
	var result string = "<html>"
	result += "<head>"
	result += "<title>ws Page</title>"
	result += "</head><body>"
	result += "<script>"
	_, proto_redirect := r.Header["X-Forwarded-Proto"]
	if proto_redirect && r.Header["X-Forwarded-Proto"][0] == "https" {
		result += "const ws = new WebSocket(\"wss://" + host + "/ws/\");"
	} else {
		result += "const ws = new WebSocket(\"ws://" + host + "/ws/\");"
	}
	result += "ws.onmessage = e => document.write(\"<div>\"+e.data+\"</div>\");"
	result += "ws.onclose = e => console.log(e.wasClean);"
	result += "</script>"
	result += "</body></html>"
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(result))

}

// две горутины на чтение и на отправку в вебсокет
// обмен через каналы
func Ws_handler(w http.ResponseWriter, r *http.Request) {
	conn, bufrw, err := acceptHandshake(w, r)

	if err != nil {
		fmt.Println(err)
		Ws_handler_start(w, r)
		return
	}
	defer conn.Close()
	chan_close := make(chan bool, 1)
	chan_write := make(chan bool) // имитация отмашки вывода в вебсокет
	go func() {
		chan_close <- false
		for { // условие окончания цикла - закрытие соединения
			f, err := readFrame(bufrw)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(f.payload)
		}
		chan_close <- true
	}()

	go func() {
		var f frame

		for {
			select {
			case msg1 := <-chan_close:
				if msg1 {
					fmt.Println("Connection closed 2")
					return
				}
			case <-chan_write:

				// отправляемое сообщение, пока жестко забитое
				f.payload = []byte{68, 117, 115, 106, 97, 32, 101, 114, 116, 101, 119, 116, 101, 119, 114, 116, 32, 208, 180, 209, 131, 209, 136, 208, 176}
				f.length = uint64(len(f.payload))
				f.isFin = true
				f.opCode = 0x1

				buf := make([]byte, 2)
				buf[0] |= f.opCode

				if f.isFin {
					buf[0] |= 0x80
				}

				if f.length < 126 {
					buf[1] |= byte(f.length)
				} else if f.length < 1<<16 {
					buf[1] |= 126
					size := make([]byte, 2)
					binary.BigEndian.PutUint16(size, uint16(f.length))
					buf = append(buf, size...)
				} else {
					buf[1] |= 127
					size := make([]byte, 8)
					binary.BigEndian.PutUint64(size, f.length)
					buf = append(buf, size...)
				}
				buf = append(buf, f.payload...)

				bufrw.Write(buf)
				bufrw.Flush()

			}
		}
	}()

	// имитация бесконечного рабочего цикла, прекращается при закрытии страницы в браузере
	for {
		time.Sleep(time.Second * 10)
		chan_write <- true
	}
}

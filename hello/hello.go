package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"

	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"work/golearn/hello/nfs/rpc"
	"work/golearn/hello/nfs/util"

	"work/golearn/hello/nfs"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/matishsiao/goInfo"
)

const Version string = "v0.2.11"

var Os string = ""

func main() {
	fmt.Println("Учебный код по Golang", Version)

	getOsParams()

	funcSyslog()
	funcInput()
	funcTypes()
	funcStructures()
	funcVariables()
	funcConstant()
	funcForIfElse()
	funcSwitch()

	funcArrays()
	funcSlices()

	funcMap()
	funcRange()
	funcClosure()

	funcInterface()

	funcErrors()
	funcGorutine()
	funcChannel()
	funcSelect()
	funcTimeout()
	funcCloseChannel()
	funcTimer()
	funcWorkerPool()
	funcAtomic()
	funcMutex()
	funcDefer()
	funcStrings()
	funcJsonToArray()

	funcTime()

	funcNumberParsing()
	func4byteToFloat()

	funcUrl()

	funcPost()

	funcDb()

	funcFileWrite()
	funcFileRead()

	funcFilePath()
	funcDir()
	funcTempFileOrDir()
	/*
				funcCommandLine()
			 funcCommandLineSubCommand()

					funcEnvironment()
					funcSpawnProcess()
					funcSignal()
					//	funcGoWithC()

		funcBase64()
		funcRandom()

			funcConnectToShare()

			funcExit(0)
			funcExit(2)
	*/
}

// Функции в Go не имеют аргументов "по умолчанию"
// Параметры ОС, на которой запущена программа
func getOsParams() {
	gi, _ := goInfo.GetInfo()
	gi.VarDump()
	Os = gi.GoOS // "windows"
}

// Типы данных
func funcTypes() {
	//Строки могут быть сложены с помощью символа +.
	fmt.Println("go" + "lang")

	//Целые числа и числа с плавающей точкой.
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)

	//Логические значения с логическими операторами
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}

// Массивы
func funcArrays() {

	fmt.Println("\nArrays and Slices")

	// В данном примере мы создаем массив a, который содержит 5 элементов с типом int.
	// Тип элементов и длина являются частью типа массива.
	// По-умолчанию массив заполняется нулевыми значениями, например, в случае int нулевое значение - 0.
	var a [5]int
	fmt.Println("emp:", a)

	//Мы можем установить значение по индексу элемента следующим образом:array[index] = value.
	// Получить значение можно аналогично - array[index].
	a[4] = 100
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])

	//Встроенная функция len возвращает длину массива.
	fmt.Println("len:", len(a))

	//Так можно инициалзировать и заполнить массив значениеми в одну строку
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Массив b: ", b)

	c := [5]string{"1", "2", "3", "4", "5"}
	fmt.Println("Массив c: ", c)

	// Тип массив является одномерным. Но вы можете совмещать типы, для создания многомерных структур.
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}

// Срезы
func funcSlices() {

	fmt.Println("\nСлайсы:")

	// В отличии от массивов, длина среза зависит от содержащихся в срезе элементов, а не определена при инициализации.
	// Создать пустой срез с ненулевой длиной можно используя оператор make.
	// В этом пример мы создаем слайс строк длиной 3 (заполненный нулевыми значениями).
	s := make([]string, 3)
	//	s := make([]string, 3, 10) // срез длиной 3 элемента и вместимостью 10 элементов
	fmt.Println("slice s: ", s)

	// Длина среза может быть переменной
	ls := 6
	sls := make([]string, ls)
	fmt.Println("slice sls: ", sls)

	// Мы можем устанавливать и получать значения, как в массивах.
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("slice s after set:", s)
	fmt.Println("get element s[2]:", s[2])

	// len возвращает длину среза, как и ожидалось.
	fmt.Println("slice length:", len(s))

	// В дополнение к базовой функциональности, срезы имеют несколько дополнительных особенностей по сравнению с массивыми.
	// Одна из них - append, которая возвращает срез, содержащий одно или более новых значений.
	// Обратите внимание, что результат функции append необходимо присвоить в переменную, т.к. это уже будет новый срез.
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("Slice after append:", s)
	// К срезу может быть добавлен другой срез
	s1 := make([]byte, 10)
	s2 := make([]byte, 4)
	s1[2] = 'A'
	s2[3] = 'B'
	s1 = append(s1, s2...)
	fmt.Println(s1)

	// Срезы могут быть скопированы с помощью copy.
	// В данном примере мы создаем пустой срез c такой же длины как и s и копируем данные из s в c.
	c := make([]string, len(s))
	copy(c, s) // copy(dst, src)
	fmt.Println("slice after copy:", c)
	// В данном примере мы создаем пустой срез cl большей длины, чем s и копируем данные из s в cl.
	// Копируется успешно, остаток заполняется нулевыми значениями для данного типа
	cl := make([]string, len(s)+10)
	copy(cl, s)
	fmt.Println("slice after copy (bigger size):", cl)
	// В данном примере мы создаем пустой срез cm меньшей длины, чем s и копируем данные из s в cm.
	// Копируется не более длины слайса-приемника
	cm := make([]string, 3)
	copy(cm, s)
	fmt.Println("slice after copy (smaller size):", cm)

	// Срезы поддерживают оператор slice (синтаксис использования slice[low:high]).
	// low входит в конечный срез, high не входит
	// Для примера, тут мы получаем срез состоящий из элементов s[2], s[3], и s[4].
	l := s[2:5]
	fmt.Println("slice[2:5] :", l)

	// Тут мы получаем срез до элемента s[5] (исключая его).
	l = s[:5]
	fmt.Println("slice[:5] :", l)

	//	А тут получаем срез от s[2] (включая его) и до конца исходного среза.
	l = s[2:]
	fmt.Println("slice2[2:] :", l)

	// Мы можем объявить и заполнить срез значениями в одну строку.
	t := []string{"go", "hren", "vam"}
	fmt.Println("slice t:", t)

	//	Срезы можно объединять в многомерные структуры данных.
	// Длина внутренних срезов может варьироваться, в отличии от многомерных массивов.
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d slice: ", twoD)
}

// Maps ассоциативный тип данных (хеши)
func funcMap() {

	fmt.Println("\nMap:")

	//Для создания пустой карты, используйте make: make(map[key_type]val_type).
	// Для map можно опускать size, будет создан map с минимальным размером
	m := make(map[string]int)

	// Вы можете установить пару ключ/значение используя привычный синтаксис map[key] = val
	m["k1"] = 7
	m["k2"] = 13

	// Вывод карты на экран с помощью fmt.Println выведет все пары ключ/значение
	fmt.Println("map m:", m)

	// Получить значение по ключу map[key].
	v1 := m["k1"]
	fmt.Println("m[k1]: ", v1)

	// Встроенная функция len возвращает количество пар ключ/значение для карты.
	fmt.Println("length :", len(m))

	// Встроенная функция delete удаляет пару key/value из карты с ключом k2.
	delete(m, "k2")
	fmt.Println("map (after delete k2):", m)

	// Необязательное второе возвращаемое значение из карты сообщает о том, существовал ли ключ в карте.
	// Это может быть использовано для устранения неоднозначности между отсутствующими ключами и ключами с нулевыми значениями, такими как 0 или “”.
	// Здесь нам не нужно само значение, поэтому мы проигнорировали его с пустым идентификатором _.
	_, key_exist := m["k2"]
	fmt.Println("k2 exists in map m:", key_exist)

	// Вы можете объявить и наполнить карту в одной строке с помощью подобного синтаксиса.
	// Мапы сортируются по ключу
	n := map[string]int{"foo": 1, "bar": 2, "aach": 3, "": 4}
	fmt.Println("map n:", n) //map n: map[:4 aach:3 bar:2 foo:1]
	nn := map[int]string{3: "foo", 1: "bar", 2: "aach", 4: ""}
	fmt.Println("map nn:", nn) //map nn: map[1:bar 2:aach 3:foo 4:]
	Actions := map[string]string{"a": "valA", "b": "valB"}
	keys := reflect.ValueOf(Actions).MapKeys()
	fmt.Println(keys) //[a b]
}

// Переменные
func funcVariables() {

	// var объявляет 1 или более переменных
	var a = "initial"
	fmt.Println(a)

	// Вы можете объявить несколько переменных за раз
	var b, c int = 1, 2
	fmt.Println(b, c)

	// Go будет определять тип по инициализированной переменной.
	var d = true
	fmt.Println(d)

	//Переменные, объявленные без соответствующей инициализации, имеют нулевое значение.
	// Например, нулевое значение для int равно 0. Для строк - пустая строка
	var e int
	fmt.Println(e)
	var es string
	fmt.Println("es=", es)

	//В Go существует короткий оператор := для объявления и инициализации переменной.
	// Например, var f string = "apple" в короткой записи превратится в
	f := "apple"
	fmt.Println(f)
}

// Константы
func funcConstant() {
	const s string = "constant"
	fmt.Println(s)

	//Оператор const может использоваться везде, где может быть использован оператор var.
	const n = 500000000

	//Постоянные выражения выполняют арифметику с произвольной точностью.
	const d = 3e20 / n
	fmt.Println(d)

	//Числовая константа не имеет типа до тех пор, пока ей не присвоен, например, при явном преобразовании.
	fmt.Println(int64(d))

	//Число может использоваться в контексте, который требует его, например, присваивание переменной или вызов функции.
	// Например, здесь math.Sin ожидает float64.
	fmt.Println(math.Sin(n))
}

// Цикл For. С помощью его же создается аналог while
// If/Else
func funcForIfElse() {

	// Стандартный тип с единственным условием (аналог while i <=3 )
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}
	//Классическая инициализация/условие/выражение после for
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// for без условия будет выполняться бесконечно пока не выполнится break (выход из цикла) или return,
	// который завершит функцию с циклом
	// аналог while(true)
	for {
		fmt.Println("loop")
		break
	}

	// Так же Вы можете использовать continue для немедленного перехода к следующей итерации цикла
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}

// if/else
// в Go не надо использовать скобки в условии, но блок необходимо заключать в фигурные скобки
// В Go нет тернарного if, поэтому вам нужно использовать полный оператор if даже для базовых условий.
// if a < b {
//	func1()
// }else{
//	func2
// }

func funcSwitch() {

	//Стандартное использование switch.
	i := 2
	fmt.Print("Write ", i, " as ")
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}

	//Вы можете использовать запятую в качестве разделителя, для перечисления нескольких значений в case.
	// Так же в данном примере используется блок по-умолчанию default.
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	default:
		fmt.Println("It's a weekday")
	}

	// switch без условия аналогичен обычному оператору if/else по своей логике.
	// Так же в этом примере что в case можно использовать не только константы.
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	// В этой конструкции switch сравниваются типы значений.
	// Вы можете использовать этот прием, для определения типа значения интерфейса.
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}
	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}
func ShowStructure(s interface{}) {
	a := reflect.ValueOf(s)
	numfield := reflect.ValueOf(s).Elem().NumField()
	if a.Kind() != reflect.Ptr {
		log.Fatal("wrong type struct")
	}
	for x := 0; x < numfield; x++ {
		fmt.Printf("Name field: `%s`  Type: `%s`\n", reflect.TypeOf(s).Elem().Field(x).Name,
			reflect.ValueOf(s).Elem().Field(x).Type())
	}
}

// Ряд (Range) перебирает элементы в различных структурах данных.
func funcRange() {

	fmt.Println("\nRange")

	// range для массивов и срезов возвращает индекс и значение для каждого элемента.
	// В данном примере мы используем range для подсчета суммы чисел в срезе.
	// Для массива синтаксис будет такой же
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)

	// Если нам не требуется индекс, мы можем использовать оператор _ для игнорирования.
	// Иногда нам действительно необходимы индексы.
	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}

	// Структуры
	type str1 struct {
		k1 string
		k2 string
	}
	type str2 struct {
		k1 int
		k2 int
		k3 str1
	}

	//	str3 := &str2{k1: 5, k2: 6, k3: str1{k1: "k1str", k2: "k2str"}}
	str3 := &str2{}
	ShowStructure(str3)

	// range для карт перебирает пары ключ/значение.
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for key, value := range kvs {
		fmt.Printf("%s -> %s\n", key, value)
	}

	//range может перебирать только ключи в карте
	for key := range kvs {
		fmt.Println("key:", key)
	}

	// range для строк перебирает кодовые точки Unicode (руны).
	// Первое значение - это начальный байтовый индекс руны, а второе - сама руна.
	// Для однобайтовых символов индекс увеличивается на 1, для двухбайтовых - на 2
	for i, c := range "go или не го" {
		fmt.Print("index=", i, " code=")
		fmt.Printf("0x%04x  %c\n", c, c)
	}
}

// Функции
// Go требует явного указания типа возвращаемого значения,
// то есть он не будет автоматически возвращать значение последнего выражения.
// Если функция принимает несколько аргументов с одинаковым типом,
// то вы можете перечислить аргументы через запятую и указать тип один раз.
// Go имеет встроенную поддержку нескольких возвращаемых значений.
// func vals(a,b, int) (int, int) {
//    return a + b, a - b
// }
//Если вы хотите получить не все значения, возвращаемые функцией,
// то можно поспользоваться пустым идентификатором _.
// Функции с переменным числом аргументов могут быть вызваны с любым количество аргументов.
// Пример такой функции - это fmt.Println.

// Замыкания (Closures)
// Функция intSeq возвращает другую функцию, которую мы анонимно определяем в теле intSeq.
// Возвращенная функция присваивается в переменную i, чтобы сформировать замыкание.
func intSeq(i int) func() int {

	return func() int {
		// Переменная i является статической для этой функции
		i++
		return i
	}
}
func funcClosure() {

	// Мы вызываем intSeq, присваивая результат (функцию) nextInt.
	// Это значение функции фиксирует свое собственное значение i,
	// которое будет обновляться каждый раз, когда мы вызываем nextInt.

	// intSeq возвращает функцию, которая возвращает инкрементированное значение числа,
	// переданного в intSeq при вызове
	nextInt := intSeq(5)
	nextInt2 := intSeq(0)

	// Посмотрите, что происходит при вызове nextInt несколько раз.
	// Чтобы подтвердить, что состояние является уникальным для этой конкретной функции, создайте и протестируйте новую.
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt2())
	fmt.Println(nextInt())
	fmt.Println(nextInt2())
}

// Структуры и методы структур
type person struct {
	name string
	age  int
}

// Функция NewPerson создает новую струкутуру person с заданным именем.
func NewPerson(name string) *person {
	// Вы можете безопасно вернуть указатель на локальную переменную,
	// так как локальная переменная переживет область действия функции.
	p := person{name: name}
	p.age = 42
	return &p
}

// Методы могут принимать получателя как указатели, так и значения.
// Важно! Если в метод передается значение, то сама структура не изменяется!
// Изменить параметры самой структуры можно только передав  ее по сссылке (через указатель)
func (pp *person) set_name() string {
	pp.name = "Хихи"
	return pp.name
}
func (p person) set_age() int {
	p.age = 13
	return p.age
}

func funcStructures() {

	fmt.Println("\nStructures")

	// Объявление струтуры как нового типа
	type Rect struct {
		width, height float64
		name          string
	}
	type Size struct {
		size uint64
	}
	// Включение одной структуры в другую
	type BigRect struct {
		Rect
		//		size: uint64 // смешанный тип не прокатывает
		Size // а так можно
	}
	var br BigRect
	br.height = 50
	br.size = 9000

	//Так создается новая структура
	fmt.Println(person{"Bob", 20})

	// Вы можете задавать имена для корректного присваивания значений при создании структуры
	fmt.Println(person{name: "Alice", age: 30})

	// Пропущенные поля будут нулевыми.
	fmt.Println(person{name: "Fred"})

	// Префикс & возвращает указатель на структуру.
	fmt.Println(&person{name: "Ann", age: 40})

	// Можно инкапсулировать создание новой структуры в функцию
	fmt.Println(NewPerson("Jon"))

	// Доступ к полям структуры осуществляется через точку.
	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)
	fmt.Println(s) // Исходная структура
	fmt.Println(s.set_age())
	fmt.Println(s) //  Передано по значению, структура не изменилась
	fmt.Println(s.set_name())
	fmt.Println(s) // Передано по указателю, структура изменилась
	fmt.Println(s.name)

	// Вы также можете использовать точки со структурными указателями - указатели автоматически разыменовываются.
	sp := &s
	fmt.Println(sp.age)

	// Структуры мутабельны.
	sp.age = 51
	fmt.Println(sp.age)
}

// Интерфейсы
// Пример базового интерфейса в Go
// В самом интерфейсе методы только объявляются, а реализуются в типах,
// которые хотят быть наследниками этого интерфейса и использоваться в аргументах функций,
// принимающих на вход интерфейс
type geometry interface {
	area() float64
	perim() float64
}

// В нашем примере мы будем реализовывать этот интерфейс для типов rect и circle.
type rect struct {
	width, height float64
	name          string
}
type circle struct {
	radius float64
	name   string
}

// Чтобы реализовать интерфейс на Go, нам просто нужно реализовать все методы в интерфейсе.
// Здесь мы реализуем интерфейс geometry для rect.
func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

// Реализация для circle.
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func (c circle) simple_func() {
	fmt.Println("simple func")
}

// Если переменная реализует интерфейс, то мы можем вызывать методы, которые находятся в этом интерфейсе.
// Функция measure использует это преимущество, для работы с любой фигурой.
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println("Площадь: ", g.area())
	fmt.Println("Периметр: ", g.perim())
}

// Используя пустой интерфейс, можем подать на вход функции любой тип
// Но сам тип при этом теряет информацию о своем исходном типе
// его надо находить при необходимости и приводить типы
func print_name(k interface{}) {
	fmt.Println("empty_interface func", k)

	switch t := k.(type) {
	case circle:
		fmt.Println(k.(circle).name)
	case rect:
		fmt.Println(k.(rect).name)
	default:
		fmt.Println("bad type", t)
	}
}
func funcInterface() {

	fmt.Println("\nInterfaces")

	r := rect{width: 3, height: 4, name: "rect"}
	c := circle{radius: 5, name: "circle"}

	// Типы circle и rect структур реализуют интерфейс geometry,
	// поэтому мы можем использовать экземпляры этих структур в качестве аргументов для measure.
	measure(r)
	measure(c)
	c.simple_func()

	var fv int = 2
	// Используя пустой интерфейс выведем имя объекта
	// Если использовать интерфейс с функцией,
	// то придется писать одинаковую функцию во всех структурах с именем
	// В таких случаях пустой интерфейс уместен
	print_name(r)
	print_name(c)
	print_name(fv)
	// понятнее, безопаснее и производительнее использовать конкретные типы,
	// то есть не пустые интерфейсные типы.

	a1 := []string{"jkjk"}
	a2 := []string{"iiuu"}
	a3 := "строка"
	if isSameType(a1, a2) {
		fmt.Println("a1 == a2 by type")
	} else {
		fmt.Println("a1 != a2 by type")
	}
	if isSameType(a1, a3) {
		fmt.Println("a1 == a3 by type")
	} else {
		fmt.Println("a1 != a3 by type")
	}

}

// проверка на совпадение типов
func isSameType(a, b interface{}) bool {
	fmt.Printf("%T\n", a)
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

// Ошибки
// Стандартная библиотека предоставляет две встроенные функции для создания ошибок:
// errors.New и fmt.Errorf.
func funcErrPrim(arg int) (int, error) {
	if arg == 43 {
		err := fmt.Errorf("error occurred at: %v", time.Now())
		return -1, err
	}
	if arg == 42 {
		//errors.New создает стандартную ошибку с указаннным сообщением
		return -1, errors.New("can't work with 42")
	}
	//Значение nil в поле ошибки, говорит о том, что ошибок нет.
	return arg + 3, nil
}
func funcErrors() {
	fmt.Println("\nErrors")

	r1, e1 := funcErrPrim(42)
	if e1 != nil {
		fmt.Println(e1)
	} else {
		fmt.Println(r1)
	}
	r2, e2 := funcErrPrim(43)
	if e2 != nil {
		fmt.Println(e2)
	} else {
		fmt.Println(r2)
	}
	r, e := funcErrPrim(44)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(r)
	}
}

// горутины
// Горутины - это легковесный тред.
func worker(id int, wg *sync.WaitGroup) {

	fmt.Printf("Worker %d starting\n", id)

	//Sleep симулирует долгую задачу.
	time.Sleep(5 * time.Second)
	fmt.Printf("Worker %d done\n", id)

	//Оповестить WaitGroup что воркер выполнился
	wg.Done()
}

func funcGorutine() {

	fmt.Println("\nGorutine")

	//Эта WaitGroup используется для ожидания выполнения всех горутин, запущенных здесь.
	var wg sync.WaitGroup

	// Запускаем несколько горутин и инкрементируем счетчик в WaitGroup для каждой запущенной горутины.
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	// анонимная функция
	wg.Add(1)
	go func(msg string) {
		time.Sleep(time.Second)
		fmt.Println(msg)
		wg.Done()
	}("going")

	//Блокируем завершение программы до момента, пока счетчик WaitGroup снова не станет равным 0.
	// Это будет означать, что все горутины выполнились.
	wg.Wait()
	fmt.Println("All goroutines ended")
}

// Каналы
// Каналы это способ связи паралелльных горутин между собой.
// Вы можете послать сообщение в канал из одной горутины и получить его в другой.
func funcChannel() {

	fmt.Print("\nChannels\n")

	// Создаем новый канал - make(chan val_type).
	// Каналы типизируются в зависимости от значений, которые они передают.
	messages := make(chan string)

	// Отправьте значение в канал, используя channel <-.
	// Здесь мы отправляем "ping" в канал messages, который мы сделали выше, из новой горутины.
	go func() {
		time.Sleep(2 * time.Second) // sleep 2 seconds
		messages <- "ping"
	}()

	// Синтаксис <-channel читает из канала.
	// Здесь мы получим сообщение "ping", которое мы отправили выше, и распечатаем его.
	msg := <-messages
	fmt.Println(msg)
	// По-умолчанию, отправление и получение блокируются, пока отправитель и получатель не будут готовы.
	// Это свойство позволило нам ждать в нашей программе сообщения "ping" без использования какой-либо другой синхронизации.

	// По умолчанию каналы не буферизованы, это означает, что они будут принимать сообщения для отправки (chan <-),
	// только если есть соответствующий канал (<- chan), готовый принять отправленное значение.
	// Буферизованные каналы принимают ограниченное количество значений без соответствующего приемника для этих значений.
	// Подобно срезам, буферизированный канал имеет длину и емкость.
	// Длина канала — это количество значений в очереди (не считанных) в буфере канала, емкость — это размер самого буфера канала.
	// Для того, чтобы вычислить длину, мы используем функцию len, а, используя функцию cap, получаем размер буфера.

	// Здесь мы создаем канал строк с буфером до 2 значений.
	buf_messages := make(chan string, 2)

	// Т.к. этот канал буферизирован, мы можем послать значения в канал без соответствующего одновременного получения.
	buf_messages <- "buffered"
	buf_messages <- "channel"

	// Позже мы можем получить эти значения как обычно.
	fmt.Print(<-buf_messages, " ")
	fmt.Println(<-buf_messages)

	// Мы можем использовать каналы для синхронизации выполнения между горутинами.
	// Вот пример использования блокирующего получения для ожидания завершения работы горутины.
	// При ожидании завершения нескольких процедур вы можете использовать WaitGroup.

	//Запустите воркера в горутине и передайте ему канал для оповещения.
	done := make(chan bool, 1)
	go worker_func(done)

	//Блокируйте, пока мы не получим уведомление от воркера из канала.
	// Без этой строки воркер может даже не успеть запуститься
	<-done
	fmt.Println("Done!")
}

// Эту функцию мы будем запускать в горутине.
// Канал done будет использован для оповещения другой горутины о том, что функция выполнена успешно.
func worker_func(done chan bool) {
	fmt.Print("working_func start...")
	time.Sleep(time.Second)
	fmt.Println("and done")

	//Отправить значение, чтобы сказать что функция выполнена успешно.
	done <- true
}

// При использовании каналов в качестве параметров функции вы можете указать,
// предназначен ли канал только для отправки или получения значений.
// Эта возможность повышает безопасность программы.
// func ping(pings chan<- string, msg string) функция ping принимает канал только для отправки значений
// func pong(pings <-chan string, pongs chan<- string) функция принимает канал для приема (pings) и канал для отправки (pongs)
// Однонаправленный канал также создается с использованием make, но с дополнительным стрелочным синтаксисом.
// roc := make(<-chan int)
// soc := make(chan<- int)

// Select
// select позволяет вам ждать нескольких операций на канале.
// select c default является неблокируемым, если нет готовых условий для case выполняется default
func Generator() chan int {
	ch := make(chan int)
	go func() {
		n := 1
		for {
			// select блокируется до тех пор, пока один из его блоков case не будет готов к запуску,
			// а затем выполняет этот блок. Если сразу несколько блоков могут быть запущены,
			// то выбирается произвольный.
			select {
			case ch <- n: // после закрытия канала этот case не будет выполняться,
				// поэтому паники от записи в закрытй канал не возникнет
				n++
			case <-ch:
				return
			}
		}
	}()
	return ch
}

func funcSelect() {

	fmt.Println("\nSelect")

	number := Generator()
	fmt.Println(<-number)
	fmt.Println(<-number)
	close(number)

	//В нашем примере мы будем выбирать между двумя каналами.
	c1 := make(chan string)
	c2 := make(chan string)

	// Каждый канал получит значение через некоторое время, например,
	// для моделирования блокировки RPC-операций, выполняемых в параллельных горутинах.
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// select выберет первый пришедший результат
	select {
	case msg1 := <-c2:
		fmt.Println("received", msg1)
	case msg2 := <-c1:
		fmt.Println("received", msg2)
	}
	fmt.Println("--")

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// Мы будем использовать select, чтобы ожидать оба значения одновременно,
	// печатая каждое из них по мере поступления.
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c2:
			fmt.Println("received", msg1)
		case msg2 := <-c1:
			fmt.Println("received", msg2)
		}
	}
	// Общее время выполнения составляет всего ~2 секунды,
	// так как и 1, и 2 секунды Sleeps выполняются одновременно.
}

// Тайм-ауты
func funcTimeout() {

	fmt.Println("\nTimeout")

	// Обратите внимание, что канал буферизован, поэтому отправка в goroutine неблокирующая.
	// Это обычная схема предотвращения утечек горутин в случае, если канал никогда не читается
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	// Вот select, реализующий тайм-аут.
	// res := <-c1 ожидает результата,
	// а <-Time.After ожидает значения, которое будет отправлено после истечения времени ожидания 1с.
	// Поскольку select берет первый полученный запрос и продолжает работу,
	// мы возьмем тайм-аут, если операция займет больше разрешенных 1с.
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	// Если мы допустим время ожидания более 3с,
	// то получение от c2 будет успешным, и мы распечатаем результат.
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}

// Стандартные отправители и получатели в каналах являются блокирующими.
// Тем не менее, мы можем использовать select с выбором по умолчанию для реализации неблокирующих отправителей,
// получателей и даже неблокирующего множественного select‘а.
/*
Мы можем использовать несколько case‘ов перед действием по-умолчанию для реализации многоцелевого неблокирующего выбора.
Здесь мы пытаемся получить неблокирующее получение messages и signals.

    select {
    case msg := <-messages:
        fmt.Println("received message", msg)
    case sig := <-signals:
        fmt.Println("received signal", sig)
    default:
        fmt.Println("no activity")
    }
*/

// Закрытие каналов
// Закрытие канала означает, что по нему больше не будет отправлено никаких значений.
// Это может быть полезно для сообщения получателям о завершении.
func funcCloseChannel() {

	fmt.Println("\nCloseChannel")

	// В этом примере мы будем использовать канал jobs для передачи задания,
	// которое должно быть выполнено из main() в горутине.
	// Когда у нас больше не будет заданий для воркера, мы закроем канал jobs.
	jobs := make(chan int, 5)
	done := make(chan bool)

	// Here’s the worker goroutine. It repeatedly receives from jobs with j, more := <-jobs.
	// In this special 2-value form of receive, the more value will be false
	// if jobs has been closed and all values in the channel have already been received.
	// We use this to notify on done when we’ve worked all our jobs.
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	// This sends 3 jobs to the worker over the jobs channel, then closes it.
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	// We await the worker using the synchronization approach we saw earlier.
	// b будет равно false, если канал закрыт
	a, b := <-done //true true
	fmt.Println(a)
	fmt.Println(b)

	a1, b1 := <-jobs // 0 false
	fmt.Println(a1)
	fmt.Println(b1)

	// Перебор значений из каналов
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// Этот range будет перебирать каждый элемент полученный из канала queue.
	// Но т.к. мы закрыли канал ранее, перебор элементов завершится после получения двух элементов.
	for elem := range queue {
		fmt.Println(elem)
	}
	// Этот пример так же демонстрирует, что возможно прочитать данные из канала уже после его закрытия.
}

// Таймер и тикер
func funcTimer() {

	fmt.Println("\nTimer")

	//	Таймеры позволяет выполнить одно событие в будущем.
	// Вы сообщаете таймеру, как долго вы хотите ждать, и он предоставляет канал,
	// который будет уведомлен в это время. Этот таймер будет ждать 2 секунды.
	timer1 := time.NewTimer(2 * time.Second)

	//<-timer1.C блокирует канал таймера C пока не будет отправлено сообщение о том, что таймер истек
	<-timer1.C
	fmt.Println("Timer 1 expired")

	// Если бы вы просто хотели подождать, вы могли бы использовать time.Sleep.
	// Одна из причин, по которой таймер может быть полезен, заключается в том,
	// что вы можете отменить таймер до его истечения, как в этом примере.
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}

	// тикеры позволяют повторять действия через определенные интервалы

	// Тикеры используют тот же механизм, что и таймеры: канал, в который посылаются значения.
	// Здесь мы будем использовать range для чтения данных из канала, которые будут поступать в него каждые 500мс.
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	// Тикеры могут быть остановлены так же как и таймеры.
	// Когда тикер будет остановлен, он не сможет больше принимать значения в свой канал.
	// Мы остановим его через 1600мс.
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

// Пул воркеров
// Это воркер, который мы будем запускать в нескольких параллельных инстансах.
// Эти воркеры будут получать задания через канал jobs и отсылать результаты в results.
// Мы будем ожидать одну секунду для каждого задания для имитации тяжелого запроса.
func worker_pool(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func funcWorkerPool() {

	fmt.Println("\nWorker Pool")

	// Чтобы использовать наш воркер пул, нам нужно отправить им задание и получить результаты выполнения.
	// Для этого мы делаем 2 канала.
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// Стартуем 3 воркерa, первоначально заблокированных, т.к. еще нет заданий.
	for w := 1; w <= 3; w++ {
		go worker_pool(w, jobs, results)
	}

	// Посылаем 6 заданий (jobs) и затем закрываем канал, сообщая о том что все задания отправлены.
	for j := 1; j <= 6; j++ {
		jobs <- j
	}
	close(jobs)

	// Наконец мы собираем все результаты. Это также гарантирует, что горутины закончились.
	// Альтернативный способ ожидания нескольких процедур заключается в использовании WaitGroup.
	for a := 1; a <= 6; a++ {
		<-results
	}
	/*
		// Проверим атомарность каналов, 5 горутин будут писать в канал
		// 5 будут читать из канала
		fmt.Println("Test channels atomic")
		var wg sync.WaitGroup
		wr := make(chan int, 100)
		for i := 1; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				j := 0
				for j < 20 {
					wr <- i
					time.Sleep(time.Second)
					j++
				}
				wg.Done()
			}(i)
		}
		for i := 1; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				j := 0
				for j < 20 {
					fmt.Println(i, <-wr)
					time.Sleep(time.Second)
					j++
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		fmt.Println("Test channels atomic is end")
	*/
}

// Атомарный счетчик
func funcAtomic() {

	fmt.Println("\nAtomic")

	// Счетчик
	var ops uint64

	// WaitGroup поможет нам подождать, пока все горутины завершат свою работу.
	var wg sync.WaitGroup

	// Мы запустим 50 горутин, каждая из которых увеличивает счетчик ровно на 1000.
	for i := 0; i < 50; i++ {
		wg.Add(1)

		// Для атомарного увеличения счетчика мы используем AddUint64,
		// присваивая ему адрес памяти нашего счетчика ops с префиксом &.
		go func() {
			for c := 0; c < 1000; c++ {
				atomic.AddUint64(&ops, 1)
			}
			wg.Done()
		}()
	}

	// Ждем пока завершатся горутины.
	wg.Wait()

	// Теперь доступ к ops безопасен, потому что мы знаем, что никакие другие горутины не пишут в него.
	// Безопасное чтение атомарного счетчика во время его обновления также возможно, используя функцию atomic.LoadUint64.
	fmt.Println("ops:", ops)
}

// Мьютексы
func funcMutex() {

	fmt.Println("\nMutex")

	var state = make(map[int]int)

	// Этот mutex будет синхронизировать доступ к state.
	var mutex = &sync.Mutex{}

	//Мы будем отслеживать, сколько операций чтения и записи мы выполняем.
	var readOps uint64
	var writeOps uint64

	// Здесь мы запускаем 100 горутин для выполнения повторных операций чтения по состоянию,
	// один раз в миллисекунду в каждой горутине.
	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {

				// Для каждого чтения мы выбираем ключ для доступа, блокируем mutex с помощью Lock() ,
				// чтобы обеспечить исключительный доступ к состоянию, читаем значение в выбранном ключе,
				// разблокируем мьютекс Unlock() и увеличиваем количество readOps.
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)

				// Немного ждем между чтениями.
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Запустим так же 10 горутин для симуляции записи, так же как мы делали для чтения.
	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Пусть 10 горутин работают над состоянием и мьютексом на секунду.
	time.Sleep(time.Second)

	// Смотрим финальное количество операций
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

	// С окончательной блокировкой состояния смотрим, как все закончилось.
	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

// Defer
// Defer используется, чтобы гарантировать, что вызов функции будет выполнен позже при выполнении программы,
// обычно для целей очистки.
func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}
func writeFile(f *os.File) {
	fmt.Println("writing")
	fmt.Fprintln(f, "data")
}

// Важно проверять наличие ошибок при закрытии файла, даже в отложенной функции.
func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// TODO: доработать  под Windows
func funcDefer() {

	fmt.Println("\nDefer")
	if Os == "windows" {
		return
	}
	// Сразу же после получения объекта файла с помощью createFile мы откладываем закрытие этого файла с помощью closeFile.
	// Она будет выполнена в конце включающей функции (main) после завершения writeFile.
	f := createFile("/tmp/defer.txt")
	defer closeFile(f)
	writeFile(f)

	// Если программа вышла по os.Exit(), то defer не будет выполнен

}

// Строковые функции
func funcStrings() {

	fmt.Println("\nString")

	var p = fmt.Println

	// Пакет strings
	p("Contains:  ", strings.Contains("test", "es"))
	p("Contains:  ", strings.Contains("тест", "ес"))
	p("Count:     ", strings.Count("test", "t"))
	p("HasPrefix: ", strings.HasPrefix("test", "te"))
	p("HasSuffix: ", strings.HasSuffix("test", "st"))
	p("Index:     ", strings.Index("test", "e"))
	p("Join:      ", strings.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", strings.Repeat("a", 5))
	p("Replace:   ", strings.Replace("foo", "o", "0", -1))
	p("Replace:   ", strings.Replace("foo", "o", "0", 1))
	p("Split:     ", strings.Split("a-b-c-d-e-б", "-"))
	p("ToLower:   ", strings.ToLower("TESTТЕСТ"))
	p("ToUpper:   ", strings.ToUpper("testтест"))
	p()
	// Примеры ниже не относятся к пакету strings, но о них стоит упомянуть -
	// это механизмы для получения длины строки и получение символа по индексу.
	p("Len: ", len("hello"))
	p("Char:", "hello"[1])
	p("Len: ", len("Привет"))
	p("Char:", "Привет"[2])

	// Форматированный вывод
	/*
		The default format for %v is:
		bool:                    %t
		int, int8 etc.:          %d
		uint, uint8 etc.:        %d, %#x if printed with %#v
		float32, complex64, etc: %g
		string:                  %s
		chan:                    %p
		pointer:                 %p
	*/
	const name, age = "Георгий", 59
	str := fmt.Sprintf("Форматирование в строку: %s is %d years old.\n", name, age)
	fmt.Println(str)

}

// Json -> Slice
// bool for JSON booleans,
// float64 for JSON numbers,
// string for JSON strings, and
// nil for JSON null.
func funcJsonToArray() {
	fmt.Println("\nJson to Slice")

	var myStoredVariable map[string]interface{}
	js_string := `{"brightness":127,"current":107,"last_seen":1656334472970,"linkquality":171,"power":23,"state":"ON","voltage":231}`
	json.Unmarshal([]byte(js_string), &myStoredVariable)

	fmt.Printf("%+v \n", myStoredVariable)
	for a, b := range myStoredVariable {
		fmt.Println(a)
		fmt.Println(reflect.TypeOf(b))
		switch t := b.(type) {
		case string:
			fmt.Println(b.(string))
		case float64:
			fmt.Println(b.(float64))
		case bool:
			fmt.Println(b.(bool))
		default:
			fmt.Printf("-%v %v \n", b, t)
		}

	}
}

// Время
type CustomDate struct {
	time.Time
}

const layout = "2006-01-02 15:04:05"

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		return
	}
	c.Time, err = time.Parse(layout, s)
	return
}
func (c CustomDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(layout))), nil
}

type Dated struct {
	DateTime CustomDate
}

func funcTime() {
	p := fmt.Println

	p("\nTime")

	///
	// Стандартные функции времени принимают ограниченный набор шаблонов
	// используем кастомный для наиболее популярного в России 2022-01-31 15:04:05
	// код не до конца ясен, еще буду разбираться (как получить таймстамп?)
	input := []byte("{\"datetime\": \"1900-01-01 12:00:04\"}")
	var d Dated
	err := json.Unmarshal(input, &d)
	if err != nil {
		p(err)
	}
	p("Unmarshal:")
	p(d.DateTime) //1900-01-01 12:00:04 +0000 UTC

	b, err := json.Marshal(d)
	if err != nil {
		p(err)
	}
	p("Marshal:")
	p(string(b)) //{"DateTime":"1900-01-01 12:00:04"}
	p("")
	////////

	// Пауза в миллисекундах
	time.Sleep(1000 * time.Millisecond)

	// Вы можете построить структуру time, указав год, месяц, день и т.д.
	// Время всегда связано с местоположением, т.е. часовым поясом.
	utc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		utc = time.Local
	}

	// Начнем с получения текущего времени
	now := time.Now()
	p(now)                                                                                        //2023-10-16 13:22:14.2225887 +0300 MSK m=+1.044650701
	fmt.Printf("Now is: %02d%02d-%02d%02d\n", now.Month(), now.Day(), now.Minute(), now.Second()) //Now is: 1016-2214

	then := time.Date(2022, 6, 9, 17, 10, 58, 651387237, utc)
	p(then) //2022-06-09 17:10:58.651387237 +0300 MSK

	// Вы можете извлечь различные компоненты значения времени.
	p(then.Year())       //2022
	p(then.Month())      //June
	p(then.Day())        //9
	p(then.Hour())       //17
	p(then.Minute())     //10
	p(then.Second())     //58
	p(then.Nanosecond()) //651387237
	p(then.Location())   //Europe/Moscow

	// Получения дня недели доступно через метод Weekday.
	p(then.Weekday()) //Thursday

	// Эти методы сравниваются два момента времени, проверяя, происходит ли первый случай до,
	// после или одновременно со вторым, соответственно.
	p(then.Before(now)) //true
	p(then.Equal(now))  //false
	p(then.After(now))  //false

	// Метод Sub возвращает Duration, интервал между двумя временами.
	diff := now.Sub(then)
	p("diff: ", diff) //diff:  11852h11m15.571201463s

	// Мы можем вычислить продолжительность.
	p("diff.Hours", diff.Hours())             //diff.Hours 11852.187658667073
	p("diff.Minutes", diff.Minutes())         //diff.Minutes 711131.2595200244
	p("diff.Seconds", diff.Seconds())         //diff.Seconds 4.2667875571201466e+07
	p("diff.Nanoseconds", diff.Nanoseconds()) //diff.Nanoseconds 42667875571201463

	// Вы можете использовать Add, чтобы продвинуть время на заданную продолжительность, или с -, чтобы переместиться назад.
	p("Then + diff:", then.Add(diff))  //Then + diff: 2023-10-16 13:22:14.2225887 +0300 MSK
	p("Then - diff:", then.Add(-diff)) //Then - diff: 2021-01-31 20:59:43.080185774 +0300 MSK

	// Используйте Unix() или UnixNano(), чтобы получить время,
	// прошедшее с начала эпохи Unix в секундах или наносекундах соответственно
	// от заданного момента времени.
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println("now: ", now) //2023-10-16 13:22:14.2225887 +0300 MSK m=+1.044650701

	// Обратите внимание, что UnixMillis не существует, поэтому, чтобы получить миллисекунды с начала эпохи Unix,
	// вам нужно будет вручную делить наносекунды.
	millis := nanos / 1000000
	fmt.Println("secs:   ", secs)   //1697451734
	fmt.Println("millis: ", millis) //1697451734222
	fmt.Println("nanos:  ", nanos)  // 1697451734222588700
	fmt.Println("        ", 1655905090419043)

	// Вы также можете конвертировать целые секунды или наносекунды Unixtime в соответствующее время.
	fmt.Println("time.Unix(secs, 0): ", time.Unix(secs, 0))   // 2022-06-09 17:10:29 +0300 MSK
	fmt.Println("time.Unix(0, nanos): ", time.Unix(0, nanos)) // 2022-06-09 17:10:29.267220286 +0300 MSK

	t := time.Unix(1655905090419043/1000000, 0)
	t2 := time.Unix(0, 1655905090419043*1000)

	fmt.Println(t.Format(time.UnixDate))  //Wed Jun 22 16:38:10 MSK 2022
	fmt.Println(t.String())               //2022-06-22 16:38:10 +0300 MSK
	fmt.Println(t2.Format(time.UnixDate)) //Wed Jun 22 16:38:10 MSK 2022
	fmt.Println(t2.String())              //2022-06-22 16:38:10.419043 +0300 MSK

	fmt.Println()
}

// Преобразование 4-х байтового среза во float-значение
// Используется в парсинге zigbee аттрибутов
func func4byteToFloat() {
	fmt.Println("\nfunc4byteToFloat")

	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(pi)

	var voltage float32
	src := []byte{0xd1, 0x7a, 0x0c, 0x45} // source value, desired value is 224.77
	buff := bytes.NewReader(src)
	err = binary.Read(buff, binary.LittleEndian, &voltage)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Printf("voltage: %0.2f \n", voltage/10)
}

// Парсинг чисел из строкового представления
func funcNumberParsing() {

	fmt.Println("\nNumberParsing")

	// С помощью ParseFloat, параметр 64 говорит о том, сколько битов точности необходимо использовать.
	f, _ := strconv.ParseFloat("1.234", 64)
	fmt.Println(f)

	var dout2 []byte = []byte("1.23213E-6\r\n")
	var sout2 string = ""
	res, err := strconv.ParseFloat(strings.Replace(string(dout2), "\r\n", "", -1), 64)
	if err == nil {
		//		sout2 = fmt.Sprintf("%f0.3", res)
		sout2 = "postgres_cpu " + strconv.FormatFloat(res, 'f', 3, 64)
	} else {
		sout2 = "postgres_cpu " + strings.Replace(string(dout2), "\r\n", "", -1)
	}
	fmt.Println(sout2)

	// Для ParseInt 0 означает вывод базы из строки. 64 необходимо, чтобы результат соответствовал 64 битам.
	// If the base argument is 0, the true base is implied by the string's prefix following
	// the sign (if present): 2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise.
	// Also, for argument base 0 only, underscore characters are permitted as defined by the Go syntax for integer literals.
	//
	// The bitSize argument specifies the integer type that the result must fit into.
	// Bit sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64.
	// If bitSize is below 0 or above 64, an error is returned.
	i, _ := strconv.ParseInt("123", 0, 64)
	fmt.Println(i)

	// ParseInt будет распознавать числа в шестнадцатеричной системе.
	// Неявно по префиксу 0x
	d, _ := strconv.ParseInt("0x1c8", 0, 64)
	fmt.Println(d)
	// Явное указание, что число в 16-ричной форме
	d2, _ := strconv.ParseInt("1c8", 16, 64)
	fmt.Println(d2)

	// ParseUint так же доступен.
	u, _ := strconv.ParseUint("789", 0, 64)
	fmt.Println(u)

	// Atoi это удобная функция для парсинга в десятеричный int.
	k, _ := strconv.Atoi("135")
	fmt.Println(k)

	// Функции парсинга возвращают ошибку в случае некорректных аргументов.
	_, e := strconv.Atoi("wat")
	fmt.Println(e)
}

// Парсинг URL
func funcUrl() {

	fmt.Println("\nUrl parsing")

	// Мы будем разбирать этот URL, который содержит схему, аутентификационные данные,
	// хост, порт, путь, параметры и фрагмент запроса.
	s := "postgres://user:pass@host.com:5432/path1/path2?kl=vudi#fru"

	// Парсим URL и убеждаемся, что нет ошибок.
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	// Получаем схему
	fmt.Println(u.Scheme)

	// User содержит всю аутентификационную информацию;
	// используйте Username и Password если надо получить конкретное поле.
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	p, _ := u.User.Password()
	fmt.Println(p)

	// Host содержит поля хост и порт, если они определены. Воспользуйтесь SplitHostPort, чтобы разделить их.
	fmt.Println(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println(host)
	fmt.Println(port)

	// Так можно получить путь и фрагмент после #.
	fmt.Println(u.Path)
	fmt.Println(u.Fragment)
	// Получаем элементы пути
	elements := strings.Split(u.Path, "/")
	fmt.Println(elements)

	// Для получения параметров в строке вида kl=vudi используйте RawQuery.
	// Вы так же можете разобрать запрос в map. Разобранный запрос в map из строк превращается в срез строк,
	// так первый элемент будет находиться по адресу [0].
	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["kl"][0])
}

// выполнение POST-запроса с параметрами и заголовком
func funcPost() (string, error) {
	const URL_EVENTS = "http://192.168.76.95:8000/events/stat?list=true"
	urlo := URL_EVENTS
	hc := http.Client{}

	params := url.Values{}

	params.Add("login", "admin")
	params.Add("password", "admin")
	encodedData := params.Encode()
	body := strings.NewReader(encodedData)

	req, _ := http.NewRequest("POST", urlo, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := hc.Do(req)
	total := ""
	if err == nil {
		headers := resp.Header
		fmt.Printf("%v", headers)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%s", string(body))
		_, ct_exists := headers["Content-Type"]
		if ct_exists && strings.Contains(headers["Content-Type"][0], "json") {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				var answer map[string]interface{}
				err = json.Unmarshal(body, &answer)
				if err == nil {
					total = answer["total"].(string)
					return "events " + total, nil
				}
			}
		}
	}

	return "", nil

	/*
		fmt.Println("\nPOST query")
		hc := http.Client{}
		url := "https://api.ng.unilight.su/v0.1/object/get/cabinet/info"
		//	вариант с json-запросом
		query := map[string]interface{}{}
		query["uid"] = 194
		query["key"] = "c48131c63d53b61293257adc5e211238"
		query["oid"] = 51031
		query_data, _ := json.Marshal(query)
		body := bytes.NewBuffer([]byte(query_data))
		// вариант с application/x-www-form-urlencoded
		/*
			params.Add("tag_id", fmt.Sprintf("%d", tag_id))
			encodedData := params.Encode()
			body := strings.NewReader(encodedData)
	*/
	/*
		req, _ := http.NewRequest("POST", url, body)
		//	вариант с json-запросом
		req.Header.Set("Content-Type", "application/json")
		// вариант с application/x-www-form-urlencoded
		// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))

		resp, err := hc.Do(req)
		if err != nil {
			fmt.Println("error Client Do")
			return "", err
		} else {
			defer resp.Body.Close()

			fmt.Println(resp.Header)
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				fmt.Println("error ReadAll")
			}

			fmt.Printf("%s\n", string(body))
			return "Ok", nil
		}
	*/
}

// Чтение файлов
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func funcFileRead() {

	fmt.Println("\nFileRead")
	if Os == "windows" {
		return
	}
	// Возможно, самая основная задача чтения файлов - это сохранение всего содержимого файла в памяти.
	// ReadFile читает файл целиком, ошибка EOF не выбрасывается. os.File не создается.
	dat, err := os.ReadFile("/tmp/dat")
	check(err)
	fmt.Print(string(dat))

	// Вам часто может потребоваться больший контроль над тем, как и какие части файла читаются.
	// Для решения этих задач начните с открытия файла, чтобы получить значение os.File.
	f, err := os.Open("/tmp/dat") // Файл открывается в режиме Только Чтение
	check(err)
	// Большую функциональность дает команда os.OpenFile
	// В примере если файла не существует, он создается(os.O_CREATE) с правами 0755
	// файл открывается на запись/чтение (os.O_RDWR)
	f1, err := os.OpenFile("/tmp/dat", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	} else {
		f1.Close()
	}
	// Пример показывает, что файл открывается не в исключительном режиме доступа

	// Прочитаем несколько байт с начала файла.
	// Будет прочитано первые 5 байт (читает не более размера переданного массива), также выведем, что фактически было прочитано.
	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))

	// Вы так же можете получить конкретное место файла с помощью Seek и выполнить Read оттуда.
	o2, err := f.Seek(6, 0) // смещение на 6 байт вперед,
	// 0 - относительно начала файла,
	// 1- относительно текущего положения,
	// 2 - относительно конца файла
	check(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes @ %d: %v\n", n2, o2, string(b2[:n2]))

	// Пакет io предоставляет некоторые функции, которые могут быть полезны для чтения файлов.
	// Например, чтение, подобное приведенному выше, может быть более надежно реализовано с помощью ReadAtLeast.
	o3, err := f.Seek(6, 0)
	check(err)
	b3 := make([]byte, 4)
	n3, err := io.ReadAtLeast(f, b3, 4)
	check(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	// Тут нет встроенной перемотки назад, но можно использовать Seek(0, 0) для этого.
	_, err = f.Seek(0, 0)
	check(err)

	// В пакете bufio реализован буферизованный ридер, который может быть полезен из-за своей эффективности
	// при большом количестве небольших операций чтения, и из-за наличия дополнительных методов чтения, которые он предоставляет.
	// Peek возвращает следующие N байт, но не двигает указатель!!!
	// Пример побайтового считывания файла с проверкой достижения конца файла
	r4 := bufio.NewReader(f)
	var b byte
	read := make([]byte, 0, 1)

	for {
		_, err := r4.Peek(1)
		if err != nil {
			fmt.Println("Peek error: ", err) // EOF
			break
		}
		b, _ = r4.ReadByte()
		read = append(read, b)
	}
	fmt.Println(string(read))

	// Пример форматированного чтения из файла
	item_data := make(map[uint16]uint64)
	prefix := "/usr/local"
	filename := prefix + "/etc/zhub4/map_addr_test.cfg"

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("OpenFile error: ", err)
	} else {
		var shortAddr uint16
		var macAddr uint64
		var r int
		var err error = nil
		for err == nil {
			r, err = fmt.Fscanf(fd, "%4x %16x\n", &shortAddr, &macAddr)
			if r > 0 {
				item_data[shortAddr] = macAddr
				fmt.Printf("%d 0x%04x 0x%016x\n", r, shortAddr, macAddr)
			}
		}
		fd.Close()
		for a, b := range item_data {
			fmt.Printf("0x%04x : 0x%016x \n", a, b)
		}
	}
	// Закройте файл, когда вы закончите использовать его (обычно закрытие с defer‘ом делается сразу после открытия).
	f.Close()

	// Пример построчного чтения файла
	fd, err = os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("OpenFile error: ", err)
	} else {
		scan := bufio.NewScanner(fd)
		// read line by line
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
		fd.Close()
	}

}

// Запись файлов
func funcFileWrite() {

	fmt.Println("\nFileWrite")
	if Os == "windows" {
		return
	}
	// В этом примере показано вот как записать строку (или только байты) в файл.
	d1 := []byte("hello\nГеоргий\n")          // строка приводится к срезу байтов
	err := os.WriteFile("/tmp/dat", d1, 0644) //filename string, data []byte, permission fs.FileMode
	check(err)

	// Для более детальной записи откройте файл для записи.
	// Существующий файл обрезается до 0, или создается новый с правами 0666, тип os.File, режим O_RDWR.
	f, err := os.Create("/tmp/dat2")
	check(err)

	// Идиоматично откладывать закрытие с помощью defer‘a сразу после открытия файла.
	defer f.Close()

	// Вы можете записать срез байт
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	// Запись строки WriteString так же доступна, равносильна записи среза байтов
	n3, err := f.WriteString("writes\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n3)

	// Выполните синхронизацию Sync для сброса записей из памяти на диск.
	f.Sync()

	// bufio предоставляет буферизованных писателей в дополнение к буферизованным читателям, которые мы видели ранее.
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered плохо\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)

	// rune (руна) - кодовая точка Unicode
	// Пример посимвольной записи Unicode символов
	n5 := 0
	var ustring string = "buffered ужасно\n"
	for _, a := range ustring {
		/*
			fmt.Printf("%#U next rune\n", a)
			runeValue, width := utf8.DecodeRuneInString(string(a)) // возвращает символ Unicode и его длину в байтах
			fmt.Printf("%#U rune has width %d\n", runeValue, width) // U+0436 'ж'
			fmt.Printf("%c rune has width %v\n", runeValue, width)  // ж
		*/
		n5s, err := w.WriteRune(rune(a))
		check(err)
		n5 = n5 + n5s
	}

	fmt.Printf("wrote %d bytes\n", n5)

	// Используйте Flush, чтобы убедиться, что все буферизованные операции были применены к основному модулю записи.
	w.Flush()

	// Пример форматированной записи в файл. Сохраняем map[uint16]uint64 построчнов файл
	devicessAddressMap := make(map[uint16]uint64)
	devicessAddressMap[0xf217] = 0x0c4314fffe17d8a8
	devicessAddressMap[0x334f] = 0x8cf681fffe0656ef
	devicessAddressMap[0x004f] = 0x00f681fffe0656ef

	prefix := "/usr/local"
	filename := prefix + "/etc/zhub4/map_addr_test.cfg"

	fd, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		for a, b := range devicessAddressMap {
			fmt.Fprintf(fd, "%04x %016x\n", a, b)
		}
		fd.Sync()
		fd.Close()
	}
}

// Пути к файлам
func funcFilePath() {

	fmt.Println("\nFile Path")

	// Join должен использоваться для создания путей в переносимом виде.
	// Он принимает любое количество аргументов и строит из них иерархический путь.
	p := filepath.Join("dir1", "dir2", "filename")
	fmt.Println("p:", p)

	// Вы должны всегда использовать Join вместо ручного объединения с помощью слешей / или \.
	// В дополнение к обеспечению переносимости, Join также нормализует пути, удаляя лишние разделители.
	fmt.Println(filepath.Join("dir1//", "filename"))
	fmt.Println(filepath.Join("dir1/../dir1", "filename"))

	// Dir и Base могут использоваться для разделения пути к каталогу и файлу.
	// В качестве альтернативы, Split вернет оба в одном вызове.
	fmt.Println("Dir(p):", filepath.Dir(p))
	fmt.Println("Base(p):", filepath.Base(p))
	spl1, spl2 := filepath.Split(p)
	fmt.Printf("Split(p): %s %s \n", spl1, spl2)

	// Можно проверить является ли путь абсолютным.
	fmt.Println(filepath.IsAbs("dir/file"))
	fmt.Println(filepath.IsAbs("/dir/file"))

	// Некоторые имена файлов имеют расширения, следующие за точкой.
	// Мы можем получить расширение из таких имен с помощью Ext.
	// Расширение приходит с точкой!
	filename := "config.json"
	ext := filepath.Ext(filename)
	fmt.Println(ext)

	filename = "config.json.duo"
	ext = filepath.Ext(filename)
	fmt.Println(ext)

	// Выводим имя файла с удалением расширения, используя strings.TrimSuffix.
	fmt.Println(strings.TrimSuffix(filename, ext))

	// Rel находит относительный путь между двумя путями base и target.
	// Возвращает ошибку, если target не может быть получен из base.
	rel, err := filepath.Rel("a/b", "a/b/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)
	rel, err = filepath.Rel("a/b", "a/c/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)
}

// Директории
func funcDir() {

	fmt.Println("\nDirectories")

	var err error

	// директория запуска программы (не по go run !!!!)
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(rootDir)

	// Создадим новую суб-директорию в текущей рабочей папке.
	// Проверим суб-директорию, чтобы убедиться, что ее еще нет
	if !isExist("./subdir") {
		err = os.Mkdir("subdir", 0755)
		check(err)
	}
	// Когда мы создаем временную директорию, хорошим тоном является удалить ее через defer.
	// os.RemoveAll удалит директорию и все, что в ней находится (по аналогии с rm -rf).
	defer os.RemoveAll("subdir")

	// Функция-помощник для создания нового пустого файла.
	createEmptyFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0644))
	}
	createEmptyFile("subdir/file1")

	// Мы можем создать иерархию из директорий, включая все родительские, с помощью MkdirAll.
	// Это является аналогом команды mkdir -p.
	err = os.MkdirAll("subdir/parent/child", 0755)
	check(err)
	createEmptyFile("subdir/parent/file2")
	createEmptyFile("subdir/parent/file3")
	createEmptyFile("subdir/parent/child/file4")

	// ReadDir перечисляет содержимое каталогов, возвращая срез объектов os.FileInfo.
	c, err := os.ReadDir("subdir/parent")
	check(err)
	fmt.Println("Listing subdir/parent")
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// Chdir позволяет изменить текущую рабочую директорию, по аналогии с cd.
	err = os.Chdir("subdir/parent/child")
	check(err)

	// Теперь мы увидим содержимое директории subdir/parent/child, когда запросим листинг текущей директории.
	c, err = os.ReadDir(".")
	check(err)
	fmt.Println("Listing subdir/parent/child")
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// Вернемся в начало
	err = os.Chdir("../../..")
	check(err)

	// Мы также можем рекурсивно обойти каталог, включая все его подкаталоги.
	// WalkDir принимает функцию обратного вызова для обработки каждого файла или каталога, которые посетили.
	fmt.Println("Visiting subdir")
	err = filepath.WalkDir("subdir", func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(" ", p, info.IsDir())
		return nil
	})
	check(err)
}

// функция проверяет существование файла/директории
func isExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Временные файлы и директории
func funcTempFileOrDir() {

	fmt.Println("\nTemp Files and Dirs")

	// Простейший способ создания временного файла - это вызов os.CreateTemp.
	// Он создаст и откроет файл для чтения и записи. Мы использовали "" в качестве первого аргумента,
	// и поэтому os.CreateTemp создаст файл в директории по-умолчанию.
	f, err := os.CreateTemp("", "sample_")
	check(err)

	// Показать имя временного файла. В ОС на основе Unix каталог, вероятно, будет /tmp.
	// Имя файла начинается с префикса, заданного в качестве второго аргумента os.CreateTemp,
	// а остальное выбирается автоматически, чтобы параллельные вызовы всегда создавали разные имена файлов.
	fmt.Println("Temp file name:", f.Name())

	// Удалите файл после того, как мы закончим.
	// Через некоторое время ОС, скорее всего, сама очистит временные файлы, но рекомендуется делать это явно.
	defer os.Remove(f.Name())

	// Мы можем записать какую-то информацию в файл.
	_, err = f.Write([]byte{1, 2, 3, 4})
	check(err)

	// Если мы намереваемся написать много временных файлов, мы можем предпочесть создать временный каталог.
	// Аргументы os.MkdirTemp совпадают с аргументами os.CreateTemp, но он возвращает имя каталога, а не открытый файл.
	dname, err := os.MkdirTemp("", "sample_dir_")
	check(err)
	fmt.Println("Temp dir name:", dname)
	defer os.RemoveAll(dname)

	// Теперь мы можем синтезировать временные имена файлов, добавив к ним префикс нашего временного каталога.
	fname := filepath.Join(dname, "file1")
	err = os.WriteFile(fname, []byte{1, 2}, 0666)
	check(err)
}

// Аргументы командной строки
// Go предоставляет пакет flag, поддерживающий базовый парсинг флагов командной строки.
func funcCommandLine() {

	fmt.Println("\nCommand line arguments")

	// os.Args предоставляет доступ к необработанным аргументам командной строки.
	// Обратите внимание, что первое значение в этом срезе - это путь к программе,
	// а os.Args [1:] содержит аргументы программы.
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	countArgs := len(os.Args[1:])

	// Вы можете получить отдельные аргументы с обычной индексацией.
	last_arg := os.Args[countArgs]
	fmt.Println("Все аргументы командной строки: ", argsWithProg)
	fmt.Println("Аргументы для программы: ", argsWithoutProg)
	fmt.Println("last argument: ", last_arg)

	// Основные объявления флагов доступны для строковых, целочисленных и логических параметров.
	// Функция flag.String возвращает строковый указатель (не строковое значение);
	// мы увидим, как использовать этот указатель ниже.
	// Флаги в командной строке дожны быть с дефисом (-word=ttt)
	// Если флаг не описан, но передан в командной строке - выбросится исключение !!!
	// Здесь мы объявляем строковой флаг word со значением по умолчанию "foo" и кратким описанием.
	wordPtr := flag.String("word", "foo", "a string")
	// Объявляем флаги numb и fork, используя тот же подход, что и выше.
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	// Также возможно вызвать метод, который использует существующую переменную,
	// объявленную в другом месте программы.
	// Обратите внимание, что в данном случае необходимо передать указатель.
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// Как только все флаги объявлены, вызовите flag.Parse(), чтобы выполнить парсинг командной строки.
	flag.Parse()

	// Здесь мы просто выведем результат парсинга и все введеные аргументы.
	// Обратите внимание, что нам нужно разыменовать указатели, например, с помощью *wordPtr,
	// чтобы получить фактические значения.
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}

// Подкоманды командной строки
func funcCommandLineSubCommand() {

	fmt.Println("\nSubcommands")

	// Мы объявляем подкоманду, используя функцию NewFlagSet,
	// и приступаем к определению новых флагов, специфичных для этой подкоманды.
	fooCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	fooEnable := fooCmd.Bool("enable", false, "enable")
	fooName := fooCmd.String("name", "", "name")

	// Для другой подкоманды мы можем определить другие флаги.
	barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	barLevel := barCmd.Int("level", 0, "level")

	// Подкоманда ожидается в качестве первого аргумента программы.
	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	// Проверяем, какая подкоманда вызвана.
	switch os.Args[1] {

	// Для каждой подкоманды мы анализируем ее собственные флаги и имеем доступ к аргументам.
	case "foo":
		fooCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'foo'")
		fmt.Println("  enable:", *fooEnable)
		fmt.Println("  name:", *fooName)
		fmt.Println("  tail:", fooCmd.Args())
	case "bar":
		barCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'bar'")
		fmt.Println("  level:", *barLevel)
		fmt.Println("  tail:", barCmd.Args())
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		//		os.Exit(1)
	}
}

// Переменные среды
func funcEnvironment() {

	fmt.Println("\nEnvironment")

	// Чтобы установить пару ключ/значение, используйте os.Setenv.
	// Чтобы получить значение для ключа, используйте os.Getenv.
	// Это вернет пустую строку, если ключ не присутствует в среде.
	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))

	// Используйте os.Environ для вывода списка всех пар ключ/значение в среде.
	// Это возвращает спез строк в формате KEY=value.
	// Вы можете использовать strings.Split, чтобы получить ключ и значение.
	// Здесь мы печатаем все ключи.
	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair)
	}

}

// Порождающий процессы
func funcSpawnProcess() {

	fmt.Println("\nSpawning Processes")

	// Мы начнем с простой команды, которая не принимает аргументов или ввода и просто печатает что-то на стандартный вывод.
	// Хелпер exec.Command создает объект для представления этого внешнего процесса.
	dateCmd := exec.Command("date")

	//.Output - это еще один хелпер, который обрабатывает общий случай запуска команды,
	// ожидает ее завершения и сбора выходных данных.
	// Если ошибок не было, dateOut будет содержать байты с информацией о дате.
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	// Далее мы рассмотрим несколько более сложный случай,
	// когда мы направляем данные во внешний процесс на его stdin и собираем результаты из его stdout.
	grepCmd := exec.Command("grep", "hello")

	// Здесь мы явно получаем каналы ввода-вывода, запускаем процесс,
	// записываем в него некоторые входные данные, читаем полученный результат и, наконец, ожидаем завершения процесса.
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := io.ReadAll(grepOut)
	grepCmd.Wait()

	// Мы опускаем проверки ошибок в приведенном выше примере,
	// но вы можете использовать обычный шаблон if err != nil для них.
	// Мы также собираем только результаты StdoutPipe, но вы можете собирать StderrPipe точно таким же образом.
	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	// Обратите внимание, что при порождении команд нам нужно предоставить
	// явно разграниченный массив команд и аргументов вместо возможности просто передать одну строку командной строки.
	// Если вы хотите создать полную команду со строкой, вы можете использовать опцию -c в bash:
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}

// Сигналы
// Сигнал SIGHUP отправляется при потере программой своего управляющего терминала.
// Сигнал SIGINT отправляется при введении пользователем в управляющем терминале символа прерывания,
// по умолчанию это ^C (Control-C).
// Сигнал SIGQUIT отправляется при введении пользователем в управляющем терминале символа выхода,
// по умолчанию это ^\ (Control-Backslash).
// SIGTERM — это общий сигнал, используемый для завершения программы.
func funcSignal() {

	fmt.Println("\nSignal")

	// Уведомление о выходе сигнала работает путем отправки значений os.Signal в канал.
	// Мы создадим канал для получения этих уведомлений
	sigs := make(chan os.Signal, 1)

	intrpt := false // забиваю по умолчанию false, чтобы отличать - выходим по сигналу или штатно

	// signal.Notify регистрирует данный канал для получения уведомлений об указанных сигналах.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Эта горутина выполняет блокировку приема сигналов.
	// Когда она получит его, то распечатает его, а затем уведомит программу, что она может быть завершена.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		intrpt = true
	}()

	// имитация выполнения программы, во время которого может придти
	// (а может не придти) сигнал прерывания работы
	fmt.Println("Имитация рабочего цикла")
	i := 0
	flag := true
	for flag {
		if intrpt {
			break
		}
		time.Sleep(time.Second)
		i++
		flag = i < 10 // имитирую возможную установку флага из других мест
	}

	if intrpt {
		// Можно сделать какие-то завершающие действия
		fmt.Println("Зачищаем перед выходом")
		fmt.Println("exiting by signal")
	} else {
		fmt.Println("exiting normal")
	}
}

// Выход Exit
// Обратите внимание, что в отличие, например, от C,
// Go не использует целочисленное возвращаемое значение из main, чтобы указать состояние выхода.
// Если вы хотите выйти с ненулевым статусом, вы должны использовать os.Exit
// Чтобы корректно работала пара Exit  и defer, вместо прямого вызова os.Exit надо вызывать
// функцию exit, которая принимает код выхода. Если он не 0, вызывать os.Exit
// В функции exit реализовать все, что что должно быть обработано при выходе
func atexit(code int) {
	// функции, обрабатываемые по defer
	fmt.Println("Функция завершения")
	if code != 0 {
		fmt.Println("Выходим со статусом ", code)
		os.Exit(code)
	}
}
func funcExit(code int) {

	fmt.Println("\nExit")

	flag := true

	// defer не будет запускаться при использовании os.Exit,
	// поэтому это напечатается только при нормальном завершении
	defer fmt.Println("Функция завершения 1")
	// Для нормального выполнения в этом случае логику переносим в единую
	// функцию завершения atexit
	defer atexit(0)

	go func() {
		i := 0
		for flag {
			fmt.Println("Некая нормальная работа в горутине", i)
			time.Sleep(time.Second)
			i++
		}
	}()

	time.Sleep(2 * time.Second)

	i := 0
	for flag {
		fmt.Println("Некая нормальная работа ", i)

		if code != 0 {
			// Имитация фатальной ошибки
			// Выход со статусом code.
			//	os.Exit(code)
			flag = false
			atexit(code)
		}
		i++
		if i > 4 {
			flag = false
		}
		time.Sleep(time.Second)
	}

	fmt.Println("Некая нормальная работа в конце")

}

/*
func funcGoWithC() {
	fmt.Println("\nGoWithC")
	_, err := C.Hello() //We ignore first result as it is a void function
	if err != nil {
		fmt.Println(err)
	}

	aC := C.int(3)
	bC := C.int(4)
	sum, err := C.sum(aC, bC)
	if err != nil {
		fmt.Println(err)
	}

	//Convert C.int result to Go int
	res := int(sum)
	fmt.Println(res)
}
*/
// Base64 кодирование/декодирование
func funcBase64() {
	fmt.Println("\nBase64 encode decode")
	data := "abc123!?$*&()'-=@~"
	// Go supports both standard and URL-compatible base64.
	// Here’s how to encode using the standard encoder.
	// The encoder requires a []byte so we convert our string to that type.

	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)
	// Decoding may return an error,
	// which you can check if you don’t already know the input to be well-formed.
	// 8AAAAAAAAAAA== F0
	// AAAAAAAAAAA= 00
	// YEkSmo2// 0x6049129a8dbf
	// YEkSmo2///8AAAAAAAAAAA== [96 73 18 154 141 191 255 255 0 0 0 0 0 0 0 0]
	//                          0x60 49 12 9a 8d bf ff ff 00 00 00 00 00 00 00 00
	sDec, _ := b64.StdEncoding.DecodeString("YEkSmo2///8AAAAAAAAAAA==")
	//	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Printf("0x%02x \n", sDec)
	fmt.Println(string(sDec))
	fmt.Println()
	//This encodes/decodes using a URL-compatible base64 format.

	uEnc := b64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}

// ввод с клавиатуры
// Выход по q
func funcInput() {
	fmt.Println("\nВвод с клавиатуры")

	go func() {

		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			fmt.Println("Введенный текст: ", text)
			if len(text) > 0 && []byte(text)[0] == 'q' {
				return
			}
		}
	}()
}
func funcRandom() {
	fmt.Println("\nRandom")
	key := rand.Intn(5)
	val := rand.Intn(100)
	fmt.Println("Random 0-5", key)
	fmt.Println("Random 0-100", val)
}

// в Windows не работает
func funcSyslog() {
	fmt.Println("\nЗапись в syslog")
	/*
		sysLog, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_LOCAL7, "gsb_tag")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(sysLog, "This is a daemon %s with gsb_tag with fmt.Printf.", "WARNING")
		sysLog.Emerg("Это emergency запись с тегом gsb_tag с syslog.Emerg.")
		sysLog.Info("Инфо запись в лог syslog.Info")
		// читаем в логе:
		// cat /var/log/local7.log | grep gsb_tag
	*/
}
func funcConnectToShare() error {
	/*
		111 TCP/UDP = RPCBIND/Portmapper ;
		635 TCP/UDP = mountd ;
		2049 TCP/UDP = nfs ;
		4045 TCP/UDP = nlockmgr (только для NFS версии 3);
		4046 TCP/UDP = status (только для NFS версии 3).
	*/
	fmt.Println("\nNFS")
	var serverIp string = "192.168.8.115"
	var share string = "/SVM02_NFS"
	util.DefaultLogger.SetDebug(true)
	util.Infof("host=%s target=%s\n", serverIp, share)

	mount, err := nfs.DialMount(serverIp)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}
	defer mount.Close()

	auth := rpc.NewAuthUnix("unknown", 0, 3)
	auth.Gids = 3
	auth.Stamp = 0xaaaaaaaa

	v, err := mount.Mount(share, auth.Auth()) // ошибка здесь
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}
	defer v.Close()

	// discover any system files such as lost+found or .snapshot
	dirs, err := ls(v, ".")
	if err != nil {
		log.Fatalf("ls: %s", err.Error())
	}
	baseDirCount := len(dirs)
	fmt.Println(baseDirCount)

	/*
		// check the length.  There should only be 1 entry in the target (aside from . and .., et al)
		if len(dirs) != 1+baseDirCount {
			log.Fatalf("expected %d dirs, got %d", 1+baseDirCount, len(dirs))
		}

		// 10 MB file
		if err = testFileRW(v, "10mb", 10*1024*1024); err != nil {
			log.Fatalf("fail")
		}

		// 7b file
		if err = testFileRW(v, "7b", 7); err != nil {
			log.Fatalf("fail")
		}

		// should return an error
		if err = v.RemoveAll("7b"); err == nil {
			log.Fatalf("expected a NOTADIR error")
		} else {
			nfserr := err.(*nfs.Error)
			if nfserr.ErrorNum != nfs.NFS3ErrNotDir {
				log.Fatalf("Wrong error")
			}
		}

		if err = v.Remove("7b"); err != nil {
			log.Fatalf("rm(7b) err: %s", err.Error())
		}

		if err = v.Remove("10mb"); err != nil {
			log.Fatalf("rm(10mb) err: %s", err.Error())
		}

		_, _, err = v.Lookup(dir)
		if err != nil {
			log.Fatalf("lookup error: %s", err.Error())
		}

		if _, err = ls(v, "."); err != nil {
			log.Fatalf("ls: %s", err.Error())
		}

		if err = v.RmDir(dir); err == nil {
			log.Fatalf("expected not empty error")
		}

		for _, fname := range []string{"/one", "/two", "/a/one", "/a/two", "/a/b/one", "/a/b/two"} {
			if err = testFileRW(v, dir+fname, 10); err != nil {
				log.Fatalf("fail")
			}
		}

		if err = v.RemoveAll(dir); err != nil {
			log.Fatalf("error removing files: %s", err.Error())
		}

		outDirs, err := ls(v, ".")
		if err != nil {
			log.Fatalf("ls: %s", err.Error())
		}

		if len(outDirs) != baseDirCount {
			log.Fatalf("directory should be empty of our created files!")
		}

		if err = mount.Unmount(); err != nil {
			log.Fatalf("unable to umount target: %v", err)
		}
	*/
	mount.Close()
	util.Infof("Completed tests")

	return nil
}

func testFileRW(v *nfs.Target, name string, filesize uint64) error {

	// create a temp file
	f, err := os.Open("/dev/urandom")
	if err != nil {
		util.Errorf("error openning random: %s", err.Error())
		return err
	}

	wr, err := v.OpenFile(name, 0777)
	if err != nil {
		util.Errorf("write fail: %s", err.Error())
		return err
	}

	// calculate the sha
	h := sha256.New()
	t := io.TeeReader(f, h)

	// Copy filesize
	n, err := io.CopyN(wr, t, int64(filesize))
	if err != nil {
		util.Errorf("error copying: n=%d, %s", n, err.Error())
		return err
	}
	expectedSum := h.Sum(nil)

	if err = wr.Close(); err != nil {
		util.Errorf("error committing: %s", err.Error())
		return err
	}

	//
	// get the file we wrote and calc the sum
	rdr, err := v.Open(name)
	if err != nil {
		util.Errorf("read error: %v", err)
		return err
	}

	h = sha256.New()
	t = io.TeeReader(rdr, h)

	_, err = ioutil.ReadAll(t)
	if err != nil {
		util.Errorf("readall error: %v", err)
		return err
	}
	actualSum := h.Sum(nil)

	if bytes.Compare(actualSum, expectedSum) != 0 {
		log.Fatalf("sums didn't match. actual=%x expected=%s", actualSum, expectedSum) //  Got=0%x expected=0%x", string(buf), testdata)
	}

	log.Printf("Sums match %x %x", actualSum, expectedSum)
	return nil
}

func ls(v *nfs.Target, path string) ([]*nfs.EntryPlus, error) {
	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		return nil, fmt.Errorf("readdir error: %s", err.Error())
	}

	util.Infof("dirs:")
	for _, dir := range dirs {
		util.Infof("\t%s\t%d:%d\t0%o", dir.FileName, dir.Attr.Attr.UID, dir.Attr.Attr.GID, dir.Attr.Attr.Mode)
	}

	return dirs, nil
}

func funcDb() {
	host := "192.168.76.95"
	port := 5432
	user := "postgres"
	password := "12345678"
	dbname := "makves"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	fmt.Println("Connected!")

	rows, err := db.Query(`SELECT count("id")  FROM "events"`)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {

		var cnt int

		err = rows.Scan(&cnt)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		fmt.Println(cnt)
	}

	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}
}

/*
Контекст пакета определяет тип контекста, который содержит крайние сроки,
сигналы отмены и другие значения в области запроса через границы API и между процессами.
Входящие запросы к серверу должны создавать контекст, а исходящие вызовы к серверам должны принимать контекст.
Цепочка вызовов функций между ними должна распространять контекст, при необходимости заменяя его производным контекстом,
созданным с помощью WithCancel, WithDeadline, WithTimeout или WithValue.
Когда Контекст отменяется, все Контексты, производные от него, также отменяются.
Функции WithCancel, WithDeadline и WithTimeout принимают контекст (родительский)
и возвращают производный контекст (дочерний) и CancelFunc.
Вызов CancelFunc отменяет дочерний и его дочерние элементы, удаляет родительскую ссылку на дочерний элемент и останавливает все связанные таймеры.
Если не удается вызвать CancelFunc, происходит утечка дочернего и его дочерних элементов до тех пор,
пока родитель не будет отменен или не сработает таймер.
Инструмент go vet проверяет, используются ли функции CancelFunc на всех путях потока управления.
Программы, использующие контексты, должны следовать этим правилам,
чтобы поддерживать согласованность интерфейсов между пакетами и
включать инструменты статического анализа для проверки распространения контекста:

Не храните контексты внутри типа структуры; вместо этого явно передайте контекст каждой функции, которая в нем нуждается.
Контекст должен быть первым параметром, обычно называемым ctx:
func DoSomething(ctx context.Context, arg Arg) error {
    // ... use ctx ...
}
Не передавайте nil Context, даже если это разрешено функцией. Передайте context.TODO, если вы не уверены, какой контекст использовать.
Используйте значения контекста только для данных в области запроса, которые передаются процессам и API,
а не для передачи необязательных параметров функциям.
Один и тот же контекст может быть передан функциям, работающим в разных горутинах;
Контексты безопасны для одновременного использования несколькими горутинами.

ctx := context.WithValue(context.Background(), "1", "one") // base context
    ctx = context.WithValue(ctx, "2", "two") //derived context

    fmt.Println(ctx.Value("1"))
    fmt.Println(ctx.Value("2"))

See https://blog.golang.org/context for example code for a server that uses Contexts.


DefaultHTTPHost = "api-iot.mcs.mail.ru"

Background returns an empty Context. It is never canceled, has no deadline,
and has no values. Background is typically used in main, init, and tests,
and as the top-level Context for incoming requests.
func Background() Context

WithCancel returns a copy of parent whose Done channel is closed as soon as
parent.Done is closed or cancel is called.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

*/

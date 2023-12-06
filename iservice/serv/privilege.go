package serv

import (
	"fmt"
	"log"
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var moduser32 = windows.NewLazySystemDLL("user32.dll")
var _GetProcessWindowStation = moduser32.NewProc("GetProcessWindowStation")
var _GetUserObjectInformation = moduser32.NewProc("GetUserObjectInformationW")

type USEROBJECTFLAGS struct {
	fInherit  bool
	fReserved bool
	dwFlags   uint32 //WSF_VISIBLE	 0x0001
}

func GetProcessWindowStation() (handle windows.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(_GetProcessWindowStation.Addr())
	handle = windows.Handle(r0)
	if handle == 0 {
		err = errnoErr(e1)
	}
	return
}

// частный случай для получения флага WSF_VISIBLE
func GetUserObjectInformation(hObj windows.Handle) (bool, error) {
	var uof USEROBJECTFLAGS = USEROBJECTFLAGS{}
	var lenO int = 12 // binary.Size(uof) возвращает 6
	var lenNeed int = 0
	r0, _, e1 := syscall.SyscallN(_GetUserObjectInformation.Addr(), uintptr(unsafe.Pointer(hObj)), uintptr(1), uintptr(unsafe.Pointer(&uof)), uintptr(lenO), uintptr(unsafe.Pointer(&lenNeed)))
	if r0 == 1 {
		return bool((uof.dwFlags & 0x1) > 0), nil
	} else {
		return false, e1
	}
}

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

// /
func SetPrivilege() bool {
	var err error

	HWINSTA, err := GetProcessWindowStation()
	if err != nil {
		if test {
			log.Printf("-1. %s ", err.Error())
		} else {
			elog.Error(1, fmt.Sprintf("-1. %s ", err.Error()))
		}
		return false
	} else {
		if test {
			log.Printf("-1. %v ", HWINSTA) //  handle to the window station
		}
		vis, err := GetUserObjectInformation(HWINSTA)
		if err == nil {
			if test {
				log.Printf("-2. %v ", vis)
			} else {
				elog.Info(1, fmt.Sprintf("-2. %v ", vis))
			}
		}
	}
	var saAttr syscall.SecurityAttributes
	saAttr.Length = uint32(reflect.TypeOf(syscall.SecurityAttributes{}).Size())
	saAttr.InheritHandle = uint32(1)
	saAttr.SecurityDescriptor = uintptr(0)

	var si syscall.StartupInfo
	si.Cb = uint32(reflect.TypeOf(syscall.SecurityAttributes{}).Size())
	si.Desktop = windows.StringToUTF16Ptr("Winsta0\\default")
	si.Flags = windows.STARTF_USESTDHANDLES
	var hToken windows.Token

	id := windows.WTSGetActiveConsoleSessionId()

	err = windows.WTSQueryUserToken(uint32(id), &hToken)
	if err != nil {
		if test {
			log.Printf("1. %s ", err.Error())
		} else {
			elog.Error(1, fmt.Sprintf("%s ", err.Error()))
		}
		return false
	}

	procHandle := windows.CurrentProcess()
	pid, err := windows.GetProcessId(procHandle)
	if err != nil {
		if test {
			log.Printf("2. %s ", err.Error())
		}
		return false
	}

	handle, err := windows.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))

	if err != nil {
		if test {
			log.Printf("3. %s ", err.Error())
		}
		return false
	}

	defer windows.CloseHandle(handle)

	var token windows.Token

	// Find process token via win32
	err = windows.OpenProcessToken(handle,
		syscall.TOKEN_ADJUST_PRIVILEGES|
			syscall.TOKEN_QUERY|
			syscall.TOKEN_DUPLICATE|
			syscall.TOKEN_ASSIGN_PRIMARY|
			syscall.TOKEN_ADJUST_SESSIONID|
			syscall.TOKEN_READ|
			syscall.TOKEN_WRITE,
		&token)

	if err != nil {
		if test {
			log.Printf("4. %s ", err.Error())
		}
		return false
	}
	privS := "SeImpersonatePrivilege"
	lpPrivName, err := syscall.UTF16PtrFromString(privS)
	if err != nil {
		if test {
			log.Printf("5. %s ", err.Error())
		}
		return false

	}
	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, lpPrivName, &luid)
	if err != nil {
		if test {
			log.Printf("6. %s ", err.Error())
		}
		return false
	}

	var tp windows.Tokenprivileges
	tp.PrivilegeCount = 1
	tp.Privileges[0].Luid = luid
	tp.Privileges[0].Attributes = 2 // SE_PRIVILEGE_ENABLED; // 2
	var returnlen uint32
	err = windows.AdjustTokenPrivileges(
		token,
		false, //disableAllPrivileges bool,
		&tp,   //newstate *Tokenprivileges,
		1024,  //len(tp.Privileges), //buflen uint32,
		nil,
		&returnlen)
	if err != nil {
		if test {
			log.Printf("7. %s ", err.Error())
		}
		return false
	}
	/*
		// Find the token user
		tokenUser, err := token.GetTokenUser()
		if err != nil {
			return false
		}

		// Close token to prevent handle leaks
		err = token.Close()
		if err != nil {
			return false
		}

		// look up domain account by sid
		account, domain, _, err := tokenUser.User.Sid.LookupAccount("localhost")
		if err != nil {
			return false
		}

		LogChan <- fmt.Sprintf("Драйвер загрузил: %s\\%s (%s)\n", domain, account, tokenUser.User.Sid.String())
	*/
	return true
}

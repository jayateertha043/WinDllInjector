package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/jayateertha043/WinDllInjector/pkg/banner"
)

const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF
	MEM_COMMIT             = 0x00001000
	MEM_RESERVE            = 0x00002000
	MEM_RESERVE_AND_COMMIT = MEM_COMMIT | MEM_RESERVE
	PAGE_READWRITE         = 0x04
	MEM_RELEASE            = 0x00008000
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	OpenProcess        = kernel32.NewProc("OpenProcess")
	VirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	VirtualFreeEx      = kernel32.NewProc("VirtualFreeEx")
	WriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	GetProcAddress     = kernel32.NewProc("GetProcAddress")
	CreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
)

func main() {

	fmt.Println(banner.Banner)

	Dll_Path := flag.String("dll", "", "Absolute Path to DLL")
	Proc_ID := flag.Uint("pid", 0, "Process ID where DLL should be injected")
	flag.Parse()

	dll_path, pid := *Dll_Path, *Proc_ID
	if dll_path != "" && pid > 0 {

		_, err := ioutil.ReadFile(dll_path)
		if err != nil || !strings.HasSuffix(dll_path, ".dll") {
			fmt.Println("Invalid Dll Path or PID range")
			os.Exit(1)
		}

		//Get Handle to Process
		HProcess, _, _ := OpenProcess.Call(ptr(PROCESS_ALL_ACCESS), ptr(true), ptr(pid))
		if HProcess == 0 {
			fmt.Println("Opening Process Handle error")
			os.Exit(1)
		}
		fmt.Printf("[+] Process handle: %v\n", HProcess)

		dll_path_len := len(dll_path) + 1

		//Allocate Virtual Memory
		virtualAddress, _, err := VirtualAllocEx.Call(HProcess, ptr(0), ptr(dll_path_len), ptr(MEM_RESERVE_AND_COMMIT), ptr(PAGE_READWRITE))
		if virtualAddress == 0 {
			fmt.Println("Virtual Alloc failed")
			fmt.Println(virtualAddress, err)
			os.Exit(1)
		}
		fmt.Printf("\n[+] Virtual Memory Allocated:%v\n", err)

		//Write Dll path to Virtual Memory
		var bytesWritten byte
		IsMemoryWritten, _, err := WriteProcessMemory.Call(HProcess, virtualAddress, ptr(dll_path), ptr(dll_path_len), uintptr(unsafe.Pointer(&bytesWritten)))
		if IsMemoryWritten == 0 {
			fmt.Println("Dll path failed to write ")
			os.Exit(1)
		}

		fmt.Printf("\n[+] Bytes Written:%v/%v\n", dll_path_len, bytesWritten)
		fmt.Printf("\n[+] Dll path written:%v\n", err)

		//Get address of LoadLibraryA function
		loadLibraryAddress, _, err := GetProcAddress.Call(kernel32.Handle(), ptr("LoadLibraryA"))
		if loadLibraryAddress == 0 {
			fmt.Println("Failed to GetProcAddress")
			os.Exit(1)
		}
		fmt.Printf("\n[+] GetProcAddress of LoadLibraryA:%v\n", err)

		//Create Remote Thread and execute dll from dll path stored in virtual memory
		var threadid int
		HThread, _, err := CreateRemoteThread.Call(HProcess, ptr(nil), ptr(0), loadLibraryAddress, virtualAddress, ptr(0), uintptr(unsafe.Pointer(&threadid)))
		if HThread == 0 {
			fmt.Println("Failed to create remote thread")
			os.Exit(1)
		}
		fmt.Printf("\n[+] Remote Handle:%v | Remote Thread ID:%v | %v\n", HThread, threadid, err)

		//Free Allocated Virtual Memory
		size := 0
		freed, _, err := VirtualFreeEx.Call(HProcess, virtualAddress, ptr(size), ptr(MEM_RELEASE))
		if freed == 0 {
			fmt.Println("Virtual Memory failed to free. ")
			os.Exit(1)
		}
		fmt.Printf("\n[+] Virtual Memory Freed:%v\n", err)
		fmt.Printf("\n[+] DLL injected successfully for the given pid\n")
		os.Exit(0)

	} else {
		fmt.Println(dll_path, pid)
		fmt.Println("Error in input")
	}

}

func ptr(val interface{}) uintptr {
	switch val.(type) {
	case byte:
		return uintptr(val.(byte))
	case bool:
		isTrue := val.(bool)
		if isTrue {
			return uintptr(1)
		}
		return uintptr(0)
	case string:
		bytePtr, _ := syscall.BytePtrFromString(val.(string))
		return uintptr(unsafe.Pointer(bytePtr))
	case int:
		return uintptr(val.(int))
	case uint:
		return uintptr(val.(uint))
	default:
		return uintptr(0)
	}
}

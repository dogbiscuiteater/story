package main

// #include <sys/ioctl.h>
// #include <string.h>
//
// void write_to_prompt(char* s)
// {
// 		int i;
//		for (i = 0; i < strlen(s); i++) {
//			ioctl(0, TIOCSTI, &s[i]);
//		}
// }
import "C"

import (

)

func main() {
	v : NewViewer()
	writeSelectionToPrompt(v.Selection)
}

func writeSelectionToPrompt(s string) {
	C.write_to_prompt(C.CString(s))
}

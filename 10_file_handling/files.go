package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// f, err := os.Open("example.txt")
	// if(err != nil) {
	// 	panic(err)
	// }

	// fileInfo, err := f.Stat()
	// if(err != nil) {
	// 	panic(err)
	// }

	// /* Files basic info */
	// fmt.Println("file name:", fileInfo.Name())
	// fmt.Println("is a folder:", fileInfo.IsDir())
	// fmt.Println("file size:", fileInfo.Size())
	// fmt.Println("file permission:", fileInfo.Mode())
	// fmt.Println("file modified at:", fileInfo.ModTime())


	/* Reading from file */
	/* 
	Note: In case of larger size file this approach is not recommended because of loading entire file in memory
	For larger files we can make use of Streams
	*/
	data, err1 := os.ReadFile("example.txt")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println(string(data))


	/* Reading a folder */
	dirInfo, err2 := os.Open(".")
	if err2 != nil {
		panic(err2)
	}

	defer dirInfo.Close()

	// In dirInfo.ReadDir(n = no. of files that needs to be displayed)
	fileInfo, err3 := dirInfo.ReadDir(-1)
	if err3 != nil {
		panic(err3)
	}

	for _, fi := range fileInfo {
		fmt.Println(fi.Name(), fi.IsDir())
	}


	/* Creating a file */
	// newFile, err4 := os.Create("exampl2.txt")
	// if err4 != nil {
	// 	panic(err4)
	// }

	// defer newFile.Close()

	// /* Now writing into the created file */
	// newFile.WriteString("hi golang ")
	// newFile.WriteString("hi there")	// it'll keep on appending the string

	// /* Another way of writing into the file - means creating array of byte type */
	// bytes := []byte("Hello Go")
	// newFile.Write(bytes)


	/* Reading & Writing to another file(streaming fashion) */
	/* First Open the source file */
	sourceFile, err5 := os.Open("example.txt")
	if err5 != nil {
		panic(err5)
	}

	defer sourceFile.Close()

	/* Now create the Open the destination file */
	destFile, err6 := os.Create("example3.txt")
	if err6 != nil {
		panic(err6)
	}

	defer destFile.Close()

	/* Now we'll actually start data transfer using streams */
	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

	/* Streaming using infinite loop */
	for {
		rb, err7 := reader.ReadByte()
		if err7 != nil {
			if err7.Error() != "EOF" {
				panic(err7)
			}

			break
		}

		err8 := writer.WriteByte(rb)
		if err8 != nil {
			panic(err8)
		}
	}

	writer.Flush()

	fmt.Println("Written to new file successfully  ")


	/* Deleting a file */
	err9 := os.Remove("exampl2.txt")
	if err9 != nil {
		panic(err9)
	}

	fmt.Println("file deleted successfully")
}
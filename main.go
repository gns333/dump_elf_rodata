package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 定义命令行参数
	inputFileFlag := flag.String("f", "", "Input binary file path")
	flag.Parse()

	// 如果没有提供输入文件路径，则打印帮助信息
	if *inputFileFlag == "" {
		fmt.Println("Usage: dump_elf_rodata -f <input_binary>")
		flag.PrintDefaults()
		return
	}

	inputPath := *inputFileFlag

	// 打开输入文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// 解析 ELF 文件
	elfFile, err := elf.NewFile(inputFile)
	if err != nil {
		fmt.Println("Error parsing ELF file:", err)
		return
	}

	// 获取 rodata 段
	rodataSection := elfFile.Section(".rodata")
	if rodataSection == nil {
		fmt.Println("Error: .rodata section not found.")
		return
	}

	// 读取 rodata 段内容
	rodataData, err := rodataSection.Data()
	if err != nil {
		fmt.Println("Error reading .rodata section:", err)
		return
	}

	// 推断输出文件路径
	outputPath := inputPath + ".rodata"

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// 写入起始虚拟地址
	startAddr := rodataSection.Addr
	err = binary.Write(outputFile, binary.LittleEndian, startAddr)
	if err != nil {
		fmt.Println("Error writing start address:", err)
		return
	}

	// 写入 rodata 段内容
	_, err = outputFile.Write(rodataData)
	if err != nil {
		fmt.Println("Error writing .rodata section data:", err)
		return
	}

	fmt.Println("Success! Output written to", outputPath)
}

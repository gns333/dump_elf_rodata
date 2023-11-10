#include <iostream>
#include <fstream>
#include <vector>
#include <cstdint>

int main(int argc, char *argv[]) {
    if (argc != 2) {
        std::cout << "Usage: ./read_output output_binary" << std::endl;
        return 1;
    }

    const char *outputPath = argv[1];

    // 打开输出文件进行读取
    std::ifstream outputFile(outputPath, std::ios::binary);
    if (!outputFile.is_open()) {
        std::cerr << "Error opening output file." << std::endl;
        return 1;
    }

    // 读取起始虚拟地址
    uint64_t startAddr;
    if (!outputFile.read(reinterpret_cast<char*>(&startAddr), sizeof(startAddr))) {
        std::cerr << "Error reading start address." << std::endl;
        return 1;
    }

    std::cout << "Start Address: 0x" << std::hex << startAddr << std::dec << std::endl;

    // 读取 rodata 段内容
    std::vector<char> rodataData;
    outputFile.seekg(0, std::ios::end);
    size_t fileSize = outputFile.tellg();
    rodataData.resize(fileSize - sizeof(startAddr));

    outputFile.seekg(sizeof(startAddr), std::ios::beg);
    if (!outputFile.read(rodataData.data(), rodataData.size())) {
        std::cerr << "Error reading .rodata section data." << std::endl;
        return 1;
    }

    // 读取指定地址的字符串
    uint64_t funcNameAddr = 0x925e800;
    uint64_t offset = funcNameAddr - startAddr;
    if (offset >= rodataData.size()) {
        std::cerr << "Error: Function address is out of range." << std::endl;
        return 1;
    }
    const char* funcName = rodataData.data() + offset;
    std::cout << "Function Name: " << funcName << std::endl;

    return 0;
}

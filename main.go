package main

import (
	"fmt"
	"math/bits"
)

func main() {
	/*基本操作*/

	// AND 与
	fmt.Println(0b0001&0b0010 == 0b0000)

	// OR 或
	fmt.Println(0b0001|0b0010 == 0b0011)

	// XOR 异或
	fmt.Println(0b0001^0b0010 == 0b0011)

	// NOT (for int8)取反,类似 ~ 操作
	var num int8 = 0b0010
	result := ^num
	fmt.Printf("binary: %08b\n", result)                 // 输出: Binary: -0000011
	fmt.Printf("unsigned binary: %08b\n", uint8(result)) // 输出: Binary: 11111101
	fmt.Println("Decimal:", result)                      // 输出: Decimal: -3

	// Bitclear (AND NOT), 清除指定位
	fmt.Println(0b0011&^0b0010 == 0b0001)

	// Left shift (<<) 左移
	fmt.Println(1<<2 == 4)

	// Right shift (>>) 右移
	fmt.Println(1>>2 == 0)

	/*bits包中的位操作函数*/

	// Count ones, 统计1的个数
	fmt.Println(bits.OnesCount8(0b00101110) == 4)

	// Count significant bits, 统计有效位数
	fmt.Println(bits.Len8(0b00101110) == 6)
	fmt.Println(bits.Len8(0b00000000) == 0)

	// Count leading zeros, 统计前导0的个数
	fmt.Println(bits.LeadingZeros8(0b00101110) == 2)
	fmt.Println(bits.LeadingZeros8(0b00000000) == 8)

	// Count trailing ones, 统计末尾0的个数
	fmt.Println(bits.TrailingZeros8(0b00101110) == 1)
	fmt.Println(bits.TrailingZeros8(0b00000000) == 8)

	// Rotate left, 左旋, 旋转 n 位
	fmt.Println(bits.RotateLeft8(0b00101110, 3) == 0b01110001)

	// Rotate right, 右旋, 旋转 n 位
	fmt.Println(bits.RotateLeft8(0b00101110, -3) == 0b11000101)

	// Reverse bits, 反转位,末位变首位
	fmt.Println(bits.Reverse8(0b00101110) == 0b01110100)

	// Reverse bytes, 反转字节, 末字节变首字节
	fmt.Println(bits.ReverseBytes16(0x00ef) == 0xef00)

	/*整数算术运算*/

	// Multiply by 2^n 乘以2的n次方: y << n == y * 2^n
	fmt.Println(3<<8 == (3 * (1 << 8)))

	// Divide by 2^n 除以2的n次方: y >> n == y / 2^n
	fmt.Println(99>>8 == (99 / (1 << 8)))

	// Check if x is even 检查 x 是不是偶数: x & 1 == 0
	fmt.Println((100&1 == 0) == (100%2 == 0))
	fmt.Println((101&1 == 0) == (101%2 == 0))

	// Check if x is a power of 2 检查 x 是不是 2 的幂: x != 0 && (x & (x - 1)) == 0
	fmt.Println((4&(4-1) == 0)) // true
	fmt.Println((5&(5-1) == 0)) // false

	// Check if a number is divisible by 2^n 检查一个数能否被 2^n 整除: ((a >> n) << n) == a
	// 这里设置n = 3, 即检查能否被8整除
	fmt.Println(((56 >> 3) << 3) == 56) // true
	fmt.Println(((57 >> 3) << 3) == 57) // false

	// Check if x and y have opposite signs 检查 x 和 y 的符号是否相反: (x ^ y) < 0
	fmt.Println((3 ^ -1) < 0)  // true
	fmt.Println((-3 ^ -1) < 0) // false

	/*单bit变换*/

	// Set the nth bit of x to 1 设第 n 位为1,n 从右从0开始: x | (1 << n)
	fmt.Println(0b1000000|(1<<3) == 0b1001000)

	// Unset the nth bit of x to 0 设第 n 位为0,n 从右从0开始: x &^ (1 << n)
	fmt.Println(0b1001100&^(1<<3) == 0b1000100)

	// Toggle the nth bit of x 翻转第 n 位,n 从右从0开始: x ^ (1 << n)
	fmt.Println(0b1001000^(1<<3) == 0b1000000)

	// Toggle all bits except the nth bit of x 翻转除了第 n 位以外的位,n 从右从0开始: ^(x ^ (1 << n))
	res := ^(0b11111111 ^ (1 << 3))
	fmt.Printf("unsigned binary: %08b\n", uint8(res))        // 输出: Binary: 00001000
	fmt.Printf("unsigned binary: %08b\n", uint8(0b11111111)) // 输出: Binary: 11111111

	// Toggle right most bit of n 翻转最右边的 n 位: x ^^ (-1<<n)
	res = 0b11111111 ^ ^(-1 << 3)
	fmt.Printf("unsigned binary: %08b\n", uint8(res))        // 输出: Binary: 11111000
	fmt.Printf("unsigned binary: %08b\n", uint8(0b11111111)) // 输出: Binary: 11111111

	// Check if the nth bit of x is set 判定第 n 位是否为1,n 从右从0开始: x & (1 << n)!= 0
	fmt.Println(0b1001000&(1<<3) != 0) // true
	fmt.Println(0b1000000&(1<<3) != 0) // false

	// Turn off rightmost set 1-bit 将最右边的 1 置为 0,其余位不变: x & (x - 1)
	fmt.Println(0b11111111&(0b11111111-1) == 0b11111110)

	// Isolate rightmost set 1-bit 保留最右边的 1 位,其余置为 0: x & (-x)
	fmt.Println(0b11111111&(-0b11111111) == 0b00000001)

	// Right propagate rightmost 1-bit 将最右边的1位的右边所有位设置为1: x | (x - 1)
	fmt.Println(0b11101000|(0b11101000-1) == 0b11101111)

	// Turn on rightmost 0-bit 将最右边的 0 位设置为1,其余不变: x | (x + 1)
	fmt.Println(0b11111000|(0b11111000+1) == 0b11111001)

	// Isolate rightmost 0-bit 保留最右边的 0 位设置为1,其余置为 0: ^x & (x + 1)
	res = ^0b11111010 & (0b11111010 + 1)
	fmt.Printf("unsigned binary: %08b\n", uint8(res)) // 输出: Binary: 00000001
	res = ^0b11111100 & (0b11111100 + 1)
	fmt.Printf("unsigned binary: %08b\n", uint8(res)) // 输出: Binary: 00000001
	res = ^0b10011111 & (0b10011111 + 1)
	fmt.Printf("unsigned binary: %08b\n", uint8(res)) // 输出: Binary: 00100000
}

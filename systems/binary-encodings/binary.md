## 1. Hexadecimal

### 1.1 Simple Conversion

9   => 0x9
136 => 0x88
247 => 0xf7

### 1.2 CSS Colors

For eight bits each of red, green, and blue, we can represent 256*256*256 =
16,777,216 colors.

### 1.3 Say hello to hellohex

For a 17-byte file, I would expect 34 hexadecimal characters to represent it.

The first five bytes in binary are:
68       65       6c       6c       6f
01101000 01100101 01101100 01101100 01101111

## 2. Integers

### 2.1 Basic Conversion

4   => 0b100
65  => 0b1000001
105 => 0b1101001
255 => 0b11111111

0b10      => 2
0b11      => 3
0b1101100 => 108
0b1010101 => 85

### 2.2 Unsigned Binary Addition

  11111111
+ 00001101
  --------
 100001100

 In decimal, this is 255 + 13 = 268. When you only have eight-bit registers, the
 most significant bit will "overflow" and you'll be left with 0b00001100, or 12.

### 2.3 Two’s Complement Conversion

127  => 0b01111111
-128 => 0b10000000
-1   => 0b11111111
1    => 0b00000001
-14  => 0b11110010

0b10000011 => -125
0b11000100 => -60

### 2.4  Addition of Two’s Complement Signed Integers

  01111111
+ 10000000
  --------
  11111111

This is 127 + -128 in decimal, so -1 as a result makes sense.

One can negate a number in two's complement by flipping all the bits and adding
one. This is because flipping the bits and adding together gives -1, as seen in
the example above, so you have to add one to make it the additive inverse where
they sum to zero.

To subtract, just negate the second argument and add.

In  8-bit two's complement, the most significant bit has value -128 (-2^7)
In 32-bit two's complement, the most significant bit has value -2,147,483,648
(-2^31)

### 2.5 Advanced: Integer Overflow Detection

If you look at the carries you can detect overflow. If there is a carry in to
the most significant bit but not another carry out, then we have overflowed
because it means we added two positive numbers and got a negative. Similarly, if
there is no carry in but a carry out, it means we've added two negatives and
gotten a positive number.


## 3. Byte ordering

### 3.1 It’s over 9000!

`xxd -g 1 9001` gives `23 29`, which is big-endian because the most significant
byte is first.

### 3.2 TCP

The hexdump of the header is:

```sh
$ xxd -e -g 1 tcpheader
00000000: af 00 bc 06 44 1e 73 68 ef f2 a0 02 81 ff 56 00  ....D.sh......V.
```

Source port:            af 00       => 44800
Destination port:       bc 06       => 48134
Sequence number:        44 1e 73 68 => 1142846312
Acknowledgement number: ef f2 a0 02 => 4025655298

The data offset is 0x8, or just 8 decimal. This is the number of 32-bit words in
the header, which is 256 bits, or 32 bytes total. Since the minimum number of
words is 5, that means there are three extra words, or 12 bytes of optional data
in the header.

## 4. Bonus: Byte ordering and integer encoding in bitmaps

The first two bytes are 0x42 0x4d, which is 66 and 77 in decimal, or "BM"
in ASCII. This gives the "Windows 3.1x, 95, NT, ... etc." variant.

Since the header is 0x7c, or 124, then it is a BITMAPV5HEADER so he width is at
byte offset 18 and the height at offset 22, both of length 4 bytes. The byte
order is little-endian so the first bytes is the least significant.

```sh
xxd -p -s 18 -l 4 image1.bmp # 18000000 => 24 px
xxd -p -s 22 -l 4 image1.bmp # 30000000 => 48 px
xxd -p -s 18 -l 4 image2.bmp # 20000000 => 32 px
xxd -p -s 22 -l 4 image2.bmp # 40000000 => 64 px
```

image1: 24x48
image2: 32x64

At offset 28 is a two-byte value giving the bits per pixel. For both images this
comes to three bytes (24 bits) per pixel.

```sh
xxd -p -s 28 -l 2 image1.bmp # 1800 => 24 bits (3 bytes)/pixel
xxd -p -s 28 -l 2 image2.bmp # 1800 => 24 bits (3 bytes)/pixel
```

The starting address of the data portion of the image is a four-byte value at
offset 10, which for both images is starting at byte 138.

```sh
xxd -p -s 10 -l 4 image1.bmp # 8a000000 => 138 bytes
xxd -p -s 10 -l 4 image2.bmp # 8a000000 => 138 bytes
```

The data portion of image1 is simply 0xffffff repeated for each pixel, which is
just a white rectangle, since it is the max value for all the colors. The data
for image2 is 0xff0000 repeated for each pixel, which is just a red rectangle,
with full intensity for red and 0 for the other colors, since the ordering of
the colors with 24 bits/pixel is BGR and the bytes repeat as 0x0000ff.

## 5. IEEE Floating Point

### 5.1 Deconstruction

0 10000100 01010100000000000000000
- -------- -----------------------
S Exponent Mantissa

sign = positive
exponent = 132 - 127 = 5
mantissa = 1 + (2^-2 + 2^-4 + 2^-6) = 1.328125

1.328125 * 2^5 = 42.5

With the largest fixed exponent, the smallest magnitude change possible is
setting the least significant bit to 1. Doing so would give:

2^-23 * 2^127 = 2^104

The smallest fixed exponent is 00000001, which with the bias is -126.
Given that smallest fixed exponent, the smallest magnitude change possible is
setting the least significant bit to 1. Doing so would give:

2^-23 * 2^-126 = 2^-149

Given these swings in magnitude, we can conclude that IEEE 754 Floating Point
values get less precise the farther away from zero they are.

### 5.2 Advanced: Float Casting

Going from the inside out:
```c
*(unsigned int *)&f // This casts the float as an unsigned int
>> 23   // Shift right 23 bits, which shoves away the mantissa
- 0x7f  // Subtracting 127 decimal, which accounts for the bias

// What remains at this point is just the exponent value itself

t = 1U << exp // Gives 2^n, where n is the result of the steps above

// Now the second statement
r = t << (t < v)

// Here we either shift to the left one place or 0, i.e. multiply by 2 or not.
// If t is already >= v, then no shift is necessary to get the next highest power
// of 2, otherwise we multiply by 2.
```



## 6. Character encodings

### 6.1 Snowman

It should take three bytes to encode the snowman because we have a code point
that is greater than 0x07ff. We also have to account for the extra bits for the
encoding itself, which for a three byte value has patten:

1110xxxx 10xxxxxx 10xxxxxx

Filling 0x2603, which is 0010 0110 0000 0011 in binary, we get:

1110 0010 1001 1000 1000 0011

And converting each nibble to hex:

1110 0010 1001 1000 1000 0011
e    2    9    8    8    3

Or 0xe29883. We can confirm with `xxd`:

```sh
$ xxd -g 1 snowman
00000000: e2 98 83                                         ...
```

### 6.2 Hello again hellohex

From 1.3 above, we had these five bytes:

01101000 01100101 01101100 01101100 01101111

In decimal they are 104 101 108 108 111 and these are the ASCII codes for:
h e l l o

We know the file is UTF-8 because the extra bytes at the end:

11110000 10011111 10011000 10000000

This matches the pattern of UTF-8 encoding with 11110xxx in the first byte and
10xxxxxx after.

Given that the multi-byte character is UTF-8 encoded, we can pull out the code
point and look it up:

11110000 10011111 10011000 10000000
     ---   ------   ------   ------
0 0001 1111 0110 0000 0000

This is hex: 0x1f600

And looking that code point up in the
[emojipedia](https://emojipedia.org/search/?q=1f600) we get the "grinning face"
😀

### 6.3 Bonus: Ding Ding Ding!

Checking out the ASCII table, the hexadecimal "bell" is 0x7. If we `echo` three
0x07 bytes out we should get the three dings. We can have `xxd` interpret the
0x7 and "print" to stdout, which should beep the terminal.

```sh
echo "0: 07 07 07" | xxd -r
```

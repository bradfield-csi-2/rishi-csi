section .text
global binary_convert
binary_convert:
  mov rax, 0           ; set result = 0
calculate:
  movzx r8, byte [rdi] ; read the next char of input into r8
  cmp   r8, 0          ; check if string is terminated
  je    done           ; jump to return if null byte

  shl rax, 1           ; shift left (i.e. multiply by 2)
  and r8,  1           ; mask the lsb of the ASCII char to yield a 0 or 1
  add rax, r8          ; add the 0 or 1 calculated above to the total

  inc rdi              ; increment to next byte of input
  jmp calculate        ; continue converting
done:
	ret

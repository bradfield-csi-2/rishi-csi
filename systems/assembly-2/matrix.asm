section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex
  imul rcx, rdx         ; offset of beginning of row
  add r8, rcx           ; add num cols to the row offset
  mov rax, [rdi + r8*4] ; get the value from memory, accounting for 4-byte ints
	ret

section .text
global pangram
pangram:
  mov rax, 0  ; the return value, defaults to false
  mov r8,  0  ; a "seen" array that keeps track if a letter was seen
check:
  movzx rcx, byte [rdi] ; read the next char of input into r8
  cmp   rcx, 0          ; check if string is terminated
  je    done            ; jump to return if null byte

  mov r10, 1         ; a placeholder that will be shifted below
  and rcx, 11111b    ; mask the low five bits of the char so 'a'=1, 'b'=2,...
  dec rcx            ; decrement rcx so 'a' shifts 0, 'b' shifts 1, ...
  shl r10, cl        ; shift left by cl (low 8 bits of rcx); 'h' gives 10000000
  or r8, r10         ; r8 is our "seen" array, so set the bit for that letter

  inc rdi            ; move to the next character in the string
  cmp r8, 0x03ffffff ; check if the lowest 26 bits of the seen array are all 1
  jne check          ; jump and keep checking if not all 1s seen
  mov rax, 1         ; otherwise, set the return value to true and return
done:
	ret

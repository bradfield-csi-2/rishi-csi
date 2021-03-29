section .text
global pangram
pangram:
  mov rax, 0  ; the return value, defaults to false
  mov r8,  0  ; a "seen" array that keeps track if a letter was seen
check:
  movzx rcx, byte [rdi] ; read the next char of input into r8
  cmp   rcx, 0          ; check if string is terminated
  je    done            ; jump to return if null byte
  inc rdi               ; move to the next character in the string

  cmp rcx, 64        ; letters start at 65 ASCII
  jng check          ; skip this char if not at least 65

  and rcx, 11111b    ; mask the low five bits of the char so 'a'=1, 'b'=2,...
  cmp rcx, 27        ; we only care about 26 letters
  jnl check          ; skip this char if it is not at most 26

  dec rcx            ; decrement rcx so 'a' shifts 0, 'b' shifts 1, ...
  mov r10, 1         ; the placeholder that will be shifted
  shl r10, cl        ; shift left by cl (low 8 bits of rcx); 'h' gives 10000000
  or r8, r10         ; r8 is our "seen" array, so set the bit for that letter

  cmp r8, 0x03ffffff ; check if the lowest 26 bits of the seen array are all 1
  jne check          ; jump and keep checking if not all 1s seen
  mov rax, 1         ; otherwise, set the return value to true and return
done:
	ret

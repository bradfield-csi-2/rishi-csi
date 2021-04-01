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

  sub rcx, 'A'       ; subtract 'A' char so 'A'=0, 'B'=1,...
                     ; chars < 'A' will be negative and thus not set a bit
  and rcx, 11111b    ; mask the low five bits for case insensitivity
  bts r8, rcx        ; set the bit in the "seen" array to the offset in rcx

  and r8, 0x03ffffff ; mask off the low 26 bits of the seen array, more significant
                     ; bits (like punctuation) are thrown away with this mask
  cmp r8, 0x03ffffff ; check if the lowest 26 bits of the seen array are all 1
  jne check          ; jump and keep checking if not all 1s seen
  mov rax, 1         ; otherwise, set the return value to true and return
done:
	ret

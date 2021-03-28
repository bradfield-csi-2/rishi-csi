section .text
global sum_to_n
sum_to_n:
  mov rax, 0    ; set total = 0
  mov rbx, 0    ; set i = 0
loop:
  add rax, rbx  ; add i to the total
  add rbx, 1    ; increment i
  cmp rbx, rdi  ; compare i to n (n is passed n via rdi)
  jng loop      ; jump to the top of the loop if i <= n
	ret

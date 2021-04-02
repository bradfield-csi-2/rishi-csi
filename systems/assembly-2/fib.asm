section .text
global fib
fib:
  push rbx      ; push rbx onto the stack to preserve its value
  mov  rbx, rdi ; load the value n into the now-safe rbx
  mov  rax, rdi ; set the return value to n
  cmp  rdi, 1   ; if n <= 1 then return
  jle done

  lea  rdi, -1[rbx] ; n - 1
  call fib          ; call fib(n-1)
  push rax          ; push the return value of fib(n-1) to the stack

  lea  rdi, -2[rbx] ; n - 2
  call fib          ; call fib(n-2)
  pop  rcx          ; pop fib(n-1)'s return value off the stack
  add  rax, rcx     ; add to the return value of fib(n-2)
done:
  pop rbx        ; pop rbx off the stack now that we're done
	ret

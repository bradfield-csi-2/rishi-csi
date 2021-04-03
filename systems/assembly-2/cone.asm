default rel

section .text
global volume
volume:
  mulss xmm0, xmm0     ; r^2
  divss xmm1, [three]  ; h/3
  mulss xmm0, xmm1     ; (r^2)h/3
  mulss xmm0, [pi]     ; pi(r^2)h/3
 	ret

section .rodata
pi:    dd 3.14159265359
three: dd 3.0

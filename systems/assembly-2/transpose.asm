section .text
global transpose
transpose:
  ; rdi: in
  ; rsi: out
  ; rdx: rows
  ; rcx: cols
  ; pos = rIndex + (nRows * cIndex)
  mov r8, 0 ; row index
  mov r9, 0 ; col index
  mov eax, 0
iterate:
  ; calculate offset of original element (A[i][j])
  mov  r11, r8            ; set r11 = rowIndex
  imul r11, rcx           ; offset to beginning of row (nCols * rowIndex)
  add  r11, r9            ; the add the colIndex to index into that row
  mov  eax, [rdi + r11*4] ; get the value from memory, accounting for 4-byte ints

  ; calculate offset of new element (A[j][i])
  mov  r11, r9            ; set r11 = colIndex
  imul r11, rdx           ; offset to beginning of col (nRows * colIndex)
  add  r11, r8            ; then add the rowIndex to index into that col
  mov [rsi + r11*4], eax  ; copy orig element to new offset in output matrix

  inc r9                  ; increment the column
  mov eax, 0              ; use this as 0 below since cmov can't have imm values

  cmp r9, rcx             ; if colIndex >= numCols
  cmovge  r9d, eax        ;   then colIndex = 0 (we finished this row)
  lea rax, [r8 + 1]       ; set aside rowIndex + 1 for next comparison
  cmp r9, 0               ; if colIndex == 0 (i.e. did we finish a row?)
  cmove r8d, eax          ;   then rowIndex = rowIndex + 1

  cmp r8, rdx            ; if rowIndex < numRows
  jl iterate             ;   then continue transposing
  ret                    ; otherwise we're done

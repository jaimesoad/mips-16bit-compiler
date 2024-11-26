print_string:
        li r10 0 # r10 = 0
        li r11 1 # r11 = 1
    ini:
        lb r0 r2 # r0 = *r2
        eq r0 r10
        bra ret
        li r3 80
        or r0 r0 r3
        li r0 0
        add r2 r2 r11
        jmp ini

ret:
    li r2 2
    add r14 r14 r2
    jr r14


print_int:
        li R6 00
        li R10 0x0A
        li R3 0x30
        li R1 01
        li R4 00
        sub R6 R6 R1
        sb R6 R4
    cic:
        div R7 R5 R10
        rem R5 R5 R10
        add R8 R5 R3
        sub R6 R6 R1
        sb R6 R8
        eq R7 R4
        bra pr
        mov R5 R7
        jmp cic
    pr:
        mov r2 r6
        mov r15 r14
        spc r14
        jmp imp
        mov r14 r15
        jmp ret

read_int:
        li r10 00
        li r1 1
        li r2 2
        li r3 0x80
        li r6 0
        li r7 0X7F
        li r9 0xa
    cic2:
        lui r0 1
        eq r0 r0
    inf1:
        bra inf1
        eq r0 r9
        bra reg
        or r0 r0 r3
        and r0 r0 r7
        sub r6 r6 r1
        sb r6 r0
        jmp cic2
    reg:
        li r3 1
        li r4 0
        li r11 0x30
    cic3:
        lb r8 r6
        add r6 r6 r1
        sub r12 r8 r11
        mul r5 r12 r3
        add r4 r4 r5
        mul r3 r3 r9
        eq r6 r10
        bra ret
        jmp cic3
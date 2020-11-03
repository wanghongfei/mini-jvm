package com.fh;

import cn.minijvm.io.Printer;

public class ExceptionCase4Test {
    public static void main(String[] args) {
        try {
            foo();
            Printer.print(10);
        } catch (RuntimeException e) {
            Printer.print(20);
        }
    }

    public static void foo() {
        throw new RuntimeException();
    }
}

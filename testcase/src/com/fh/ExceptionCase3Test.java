package com.fh;

import cn.minijvm.io.Printer;

public class ExceptionCase3Test {
    public static void main(String[] args) {
        Printer.print(10);
        throw new RuntimeException();
    }
}

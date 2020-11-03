package com.fh;

import cn.minijvm.io.Printer;

public class ExceptionCase2Test {
    public static void main(String[] args) {
        try {
            Printer.print(10);
            throw new RuntimeException();

            //} catch (SimpleException e) {
            //    Printer.print(20);
        } finally {
            Printer.print(30);
        }
    }
}

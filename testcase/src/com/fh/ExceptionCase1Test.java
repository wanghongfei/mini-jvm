package com.fh;

import cn.minijvm.io.Printer;

public class ExceptionCase1Test {
    public static void main(String[] args) {
        try {
            Printer.print(10);
            throw new SimpleException();

        } catch (SimpleException e) {
            Printer.print(20);

        } finally {
            Printer.print(30);
        }
    }
}

package com.fh;
import cn.minijvm.io.Printer;

public class StaticFieldTest {
    private static int age = 100;
    private static int id = 200;

    public static void main(String[] args) {
        Printer.print(age);
        id = 400;
        Printer.print(id);
    }

}
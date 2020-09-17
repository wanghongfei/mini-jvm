
package com.fh;
import cn.minijvm.io.Printer;

public class ReflectionTest {
    public static void main(String[] args) throws Exception {
        Object obj = new Object();
        Class clazz = Object.class;
        Printer.printString(clazz.getName());

        clazz = obj.getClass();
        Printer.printString(clazz.getName());

        Printer.printBool(true);
        Printer.printBool(false);

        Printer.printString(clazz.toString());

//        StringBuilder sb = new StringBuilder();
//        sb.append("test");
//        Printer.printString(sb.toString());
    }
}
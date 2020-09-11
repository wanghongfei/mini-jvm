
package com.fh;
import cn.minijvm.io.Printer;

public class ReflectionTest {
    public static void main(String[] args) throws Exception {
        Object obj = new Object();
        Class clazz = Object.class;
        clazz = obj.getClass();


        Printer.printString(clazz.getName());
    }
}
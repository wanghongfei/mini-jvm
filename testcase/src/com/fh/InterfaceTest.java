
package com.fh;
import cn.minijvm.io.Printer;

public class InterfaceTest {
    public static void main(String[] args) {
        CanSay canSay = new Dog();
        canSay.say();
        canSay.say();
        Printer.print(500);
    }

    public static interface CanSay {
        void say();
    }

    public static class Dog implements CanSay {
        public void say() {
            Printer.print(100);
        }
    }
}

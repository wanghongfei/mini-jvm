
package com.fh;

public class InterfaceTest {
    public static void main(String[] args) {
        CanSay canSay = new Dog();
        canSay.say();
        canSay.say();
        print(500);
    }

    public static native void print(int num);

    public static interface CanSay {
        void say();
    }

    public static class Dog implements CanSay {
        public void say() {
            print(100);
        }
    }
}

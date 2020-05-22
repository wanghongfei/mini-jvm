# Mini-JVM

使用Go实现的Java虚拟机(不完整)，解释执行，仅学习JVM使用；

支持整数加法、循环、控制台输出、简单对象创建(不会调用构造方法)、部分继承特性、形参全部为int类型的方法调用、方法重载；

其他特性还未实现，如异常、JDK类库加载、线程等等；



## 使用方法

```shell
./mini-jvm [主类全限定性名] [classpath(仅支持.class文件所在目录,不支持jar)]
```



## 已支持的字节码

解释器已经支持的字节码如下：

```go
const (
	Nop byte = 0x00

	Iconst0 = 0x03
	Iconst1 = 0x04
	Iconst2 = 0x05

	Istore0 = 0x3b
	Istore1 = 0x3c
	Istore2 = 0x3d
	Istore3 = 0x3e

	Bipush = 0x10
	Sipush = 0x11

	Iload0 = 0x1a
	Iload1 = 0x1b
	Iload2 = 0x1c
	Iload3 = 0x1d

	Aload0 = 0x2a
	Aload1 = 0x2b
	Aload2 = 0x2c

	Astore0 = 0x4b
	Astore1 = 0x4c
	Astore2 = 0x4d

	Dup = 0x59

	Iadd = 0x60

	Iinc = 0x84

	Ificmpgt = 0xa3
	Ificmple = 0xa4
	Goto = 0xa7

	Return = 0xb1

	GetField = 0xb4
	Putfield = 0xb5

	Invokevirtual = 0xb6
	Invokespecial = 0xb7
	Invokestatic = 0xb8

	New = 0xbb

	Ireturn = 0xac
)
```



## 已实现的特性举例

`mini_jvm_test.go`中有所有用例的单元测试。



计算从1加到100的值(`testclass/com/fh/ForLoopPrintTest.java`)：

```java
package com.fh;
public class ForLoopPrintTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }
        print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static native void print(int num);
}

```

```shell
./mini-jvm ForLoopPrintTest ../testclass/
```





递归调用方法(`testclass/com/fh/RecursionTest`)：

```java
package com.fh;

public class RecursionTest {
    public static void main(String[] args) {
        // 从1打印到100
        foo(1);
    }

    public static void foo(int i) {
        if (i > 100) {
            return;
        }

        print(i);
        i++;
        foo(i);
    }

    public static native void print(int num);
}
```

```
./mini-jvm RecursionTest ../testclass/
```





简单对象的创建/访问(`testclass/com/fh/NewSimpleObjectTest.java`)：

```java
package com.fh;
public class NewSimpleObjectTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }

        Person p = new Person();
        p.setAge(sum);
        int age = p.getAge();

        print(age);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static native void print(int num);
}

```

```
./mini-jvm NewSimpleObjectTest ../testclass/
```





方法重写(`testclass/com/fh/MethodOverrideTest.java`)：

```java
package com.fh;
public class MethodOverrideTest {
    public static void main(String[] args) {
        int sum = add(100, 200);
        print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static int add(int x, int y, int z) {
        return x + y + z;
    }

    public static native void print(int num);
}

```

```
./mini-jvm MethodOverrideTest ../testclass/
```


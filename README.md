# Mini-JVM

使用Go实现的Java虚拟机(不完整)，解释执行，仅学习JVM使用；

Mini-JVM首先会从`classpath`中加载主类的class文件，然后找到main方法的字节码解释执行；执行过程中如果遇到新的类符号引用，则会通过全限定性名再从`classpath`中加载新的类文件，以此类推；

当前仅支持整数加法、循环、控制台输出、简单对象创建(不会调用构造方法)、部分继承特性、形参全部为int类型的方法调用、方法重载、方法重写、接口方法调用；其他特性还未实现，如异常、JDK类库加载、线程等等；

  

## 使用方法

编译(打开mod支持)：

```
go build -o mini-jvm
```

运行：

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
	Iconst3 = 0x06
	Iconst4 = 0x07
	Iconst5 = 0x08

	Iaload = 0x2e

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
	Iastore = 0x4f

	Dup = 0x59

	Iadd = 0x60

	Iinc = 0x84

	Ifeq = 0x99
	Ifne = 0x9a
	Iflt = 0x9b
	Ifge = 0x9c
	Ifgt = 0x9d
	Ifle = 0x9e

	Ificmpeq = 0x9f
	Ificmpne = 0xa0
	Ificmplt = 0xa1
	Ificmpge = 0xa2
	Ificmpgt = 0xa3
	Ificmple = 0xa4
	Goto = 0xa7

	Return = 0xb1

	GetField = 0xb4
	Putfield = 0xb5

	Newarray = 0xbc

	Invokevirtual = 0xb6
	Invokespecial = 0xb7
	Invokestatic = 0xb8

	New = 0xbb

	Ireturn = 0xac

	Wide = 0xc4
)
```



## 已实现的特性举例

`mini_jvm_test.go`中有所有用例的单元测试；

以下Java代码均使用Java8进行编译；

  

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

   

方法重载(`testclass/com/fh/MethodReloadTest.java`)：

```java
package com.fh;
public class MethodReloadTest {
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
./mini-jvm MethodReloadTest ../testclass/
```

  

方法重写(`testclass/com/ch/ClassExtendTest.java`)：

```java
package com.fh;

public class ClassExtendTest {
    public static void main(String[] args) {
        Person person = new Person();
        print(person.say());

        person = new Student(); // Student继承了Person
        print(person.say());
    }

    public static native void print(int num);
}
```

```
./mini-jvm ClassExtendTest ../testclass/
```



判断语句、自增语句：

```java
package com.fh;

public class IfTest {
    public static void main(String[] args) {
        int sum = 0;

        if (sum > 0) {
            sum = 100;
        }
        if (sum >= 0) {
            sum -= 200;
        }
        if (sum < 0) {
            sum -= 100;
        }
        if (sum <= 0) {
            sum -= 1;
        }
        if (sum == 0) {
            sum += 2000;
        }

        print(sum); // -301

        if (sum > 10) {

        }
        if (sum >= 10) {

        }
        if (sum < 10) {

        }
        if (sum <= 10) {

        }
        if (sum == 10) {

        }
        print(sum);
    }

    public static native void print(int num);

```

  

int数组的读写： `testclass/com/fh/ArrayTest.java`
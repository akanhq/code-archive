/*
    在函数定义中使用泛型
        需要指明T的限制条件（具备什么能力）
*/

fn largest_i32(list: &[i32]) -> &i32 {
    let mut largest = &list[0];
    for item in list {
        if item > largest {
            largest = item
        }
    }
    largest
}

fn largest_char(list: &[char]) -> &char {
    let mut largest = &list[0];
    for item in list {
        if item > largest {
            largest = item
        }
    }
    largest
}

fn largest<T: PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];
    for item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

fn main() {
    let number_list = vec![1, 2, 3, 4, 5, 6];
    let result = largest(&number_list);
    println!("{}", result);

    let char_list = vec!['a', 'b', 'c', 'd', 'e', 'f'];
    let result = largest(&char_list);
    println!("{}", result);
}

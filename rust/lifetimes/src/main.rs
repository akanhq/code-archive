use std::result;

fn main() {
    // let r;
    // {
    //     let x = 32;
    //     r = &x;
    // }
    // println!("r is {}", r);

    let string1 = String::from("long string is long");
    let result;
    {
        let string2 = "short";
        result = longest(string1.as_str(), string2);
    }
    println!("The longest string is {}", result);
}

fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() { x } else { y }
}

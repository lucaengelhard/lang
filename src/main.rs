use crate::lexer::Lexer;

mod lexer;

fn main() {
    let file_contents = match lexer::get_source_file("test/main.lang") {
        Ok(v) => v,
        Err(_e) => panic!()
    };

    println!("{}", file_contents);
    let mut lex = Lexer::new(file_contents);
    println!("{:?}", lex.parse());
}

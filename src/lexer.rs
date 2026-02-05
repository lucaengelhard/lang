use std::fs;

#[derive(Debug)]
pub enum TokenType {
    EOF,
    ILLEGAL,

    INT,
    FLOAT,
    STRING,

    PLUS,
    MINUS,
    ASTERISK,
    SLASH,
    PERCENT,
    POWER,

    SEMICOLON,
    LPAREN,
    RPAREN
}

#[derive(Debug)]
pub enum TokenLiteral {
    Str(String),
    Bool(bool),
    Int(i64),
    Float(f64)
}


#[derive(Debug)]
pub struct Token {
    kind: TokenType,
    literal: TokenLiteral,
    line_no: u64,
    char_index: u64,
}

pub struct Lexer {
    source: String,
    position: usize,
    line_no: usize
}

impl Lexer {
    pub fn new(source: impl Into<String>) -> Self {
        Self {
            source: source.into(),
            position: 0,
            line_no: 0
        }
    }

    fn current_char(&mut self) -> Option<char> {
        self.source.chars().nth(self.position)
    }

    fn next_char(&mut self) -> Option<char> {
        self.source.chars().nth(self.position + 1)
    }

    fn advance(&mut self)-> Result<char, ()> {
        if self.position >= self.source.len() {
            return Err(());
        }

        let res = match self.current_char() {
            Some(c) => Ok(c),
            None => Err(())
        };

        self.position += 1;

        res
    }

    pub fn skip_whitespace(&mut self) {
        while let Ok(char) = self.advance() {
            if char.is_whitespace() {
                if char == '\n' {
                    self.line_no += 1;
                }
                if char == '\r' {
                    self.line_no += 1;
                    if self.next_char().is_some_and(|c| c == '\n') {
                        let _ = self.advance();
                    }
                }
            } else {
                break;
            }
        }
    }
}

pub fn get_source_file(path: &str) -> Result<String, std::io::Error> {
    fs::read_to_string(path)
}
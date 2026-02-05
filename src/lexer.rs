use std::fs;

#[derive(Debug)]
pub enum Token {
    PLUS(char),

    IDENTIFIER(String),

    ILLEGAL(char),
}

#[derive(Debug)]
pub struct Lexer {
    source: String,
    pub position: usize,
}

impl Lexer {
    pub fn new(source: String) -> Self {
        Self {
            source: source.chars().collect(),
            position: 0,
        }
    }

    pub fn parse(&mut self) -> Vec<Token> {
        let mut res = Vec::new();
        while let Some(tok) = self.next_token()  {
            res.push(tok);
        };
        res
    }   

    fn current_char(&self) -> Option<char> {
        self.source.chars().nth(self.position)
    }

    fn increment_position(&mut self) -> Option<char> {
        let char = self.current_char();
        self.position += 1;
        char
    }

    fn next_token(&mut self) -> Option<Token> {
        let current = match self.current_char() {
            Some(c) => c,
            None => todo!(),
        };

        let tok: Option<Token> = match current {
            '+' => Some(Token::PLUS(current)),
            c => {
                if c.is_whitespace() {
                    None
                } else {
                    Some(self.read_identifier())
                }
            }
        };

        self.increment_position();
        tok
    }

    fn read_identifier(&mut self) -> Token {
        let start_position = self.position;
        while let Some(c) = self.increment_position()
            && c.is_alphabetic()
        {
            println!("{}", c)
        }
        println!("{start_position} {}",self.position);
        //println!("{}", self.current_char().unwrap());

        return Token::IDENTIFIER("hi".to_string());
    }
}

pub fn get_source_file(path: &str) -> Result<String, std::io::Error> {
    fs::read_to_string(path)
}

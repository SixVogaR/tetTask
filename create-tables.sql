DROP TABLE IF EXISTS currency;

CREATE TABLE currency (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(3) NOT NULL,
    value DECIMAL(10,5) NOT NULL,
    date date NOT NULL,
    PRIMARY KEY (`id`)
);
use learninghub;

CREATE TABLE User (
		userID int not null AUTO_INCREMENT primary key,
        username varchar(25) not null,
        password varchar(25)
        );

CREATE TABLE Flashcard_instance (
	flashcardId int not null AUTO_INCREMENT primary key,
    deckId int,
    term varchar(25) not null,
    definition varchar(255) not null,
	userID int,
    foreign key (userID) references User(userID)
);
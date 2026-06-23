CREATE TABLE dailyhealthupdatecache(
                                       kunstwerkid SERIAL PRIMARY KEY REFERENCES kunstwerk(id),
                                       status TEXT NOT NULL,
                                       aantalsensoren INT NOT NULL,
                                       aantalactievesensoren INT NOT NULL,
                                       aantalafwijkendesensoren INT NOT NULL,
                                       aantalafwijkingen INT NOT NULL
);
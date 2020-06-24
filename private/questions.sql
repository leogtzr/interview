-- === Java questions ===
insert into question (question, answer, topic_id, level_id) values(
"Cuál es la diferencia entre una clase abstracta y una interface?",
"", 
1, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"Do you extend from an abstract class or implement?",
"", 
1, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"What is an immutable class and how to create it?",
"", 
1, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What are the benefits of immutable classes?",
"", 
1, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What is the difference between void x(List<A> as) and void x(ArrayList<A> as)?",
"", 
1, 
2
);


insert into question (question, answer, topic_id, level_id) values(
"Why do they say that inheritance is broken in Java?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What is the risk of declaring constants within an interface? e.g interface X { int VALUE = 3; }",
"", 
1, 
3
);

-- some Collections framework related questions
/*
insert into question (question, answer, topic_id, level_id) values(
"What is fail-fast in Java Iterators?",
"Iterator fail-fast property checks for any modification in the structure of the underlying collection everytime we try to get the next element. \n
If there are any modifications found, it throws ConcurrentModificationException", 
1, 
3
);
*/

insert into question (question, answer, topic_id, level_id) values(
"What is the difference between fail-fast VS fail-safe?",
"", 
1, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"When you create an ArrayList, you can specify an initial capacity in the constructor, what does it mean?",
"", 
1, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What is a 'load factor' in a HashSet<T>?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the collection views in the Map interface?",
"Set<K> keySet(), Collection<V> values, Set<Map.Entry<K, V>> entrySet()", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the collection views in the Map interface?",
"Set<K> keySet(), Collection<V> values, Set<Map.Entry<K, V>> entrySet()", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the differences between a HashMap and a Hashtable?",
"HMap allows null keys and null values. Htable doesn't allow null keys and null values. Htable is synchronized. Htable is kind of legacy.", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the differences between a TreeMap and a HashMap?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the differences between a Vector and an ArrayList?",
"Vector is synchronized and the way they grow up is different.", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"Who is faster? an ArrayList or a Vector and why?",
"", 
1, 
3
);


insert into question (question, answer, topic_id, level_id) values(
"How can you create a synchronized collection from given collection?",
"Collections.synchronizedCollection(Collection c)", 
1, 
3
);

-- Object Oriented Programming
insert into question (question, answer, topic_id, level_id) values(
"Can you change the return type when you override a method?",
"Yes, covariant return.", 
1, 
3
);

-- Other
insert into question (question, answer, topic_id, level_id) values(
"What is the difference between a static and a default method in an interface in Java 8?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What is a marker interface?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What can you tell me about the following class declaration? public final abstract class Animal{}",
"", 
1, 
3
);

-- S.O.L.I.D
insert into question (question, answer, topic_id, level_id) values(
"What is the Liskov Substituion Principle?",
"
It means that the classes fellow developer created by extending our class should be able to 
fit in application without failure. This requires the objects of your subclasses to behave 
in the same way as the objects of your superclass. 
", 
1, 
3
);


insert into question (question, answer, topic_id, level_id) values(
"What is the Interface Segregation Principle?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What do you understand by 'Depend on abstractions, not concretions'?",
"", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What are the advantages of a static factory method VS a constructor?",
"Readable Names, Polymorphism, Coupling (too many new X...), caching, etc", 
1, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"Why would you prefer composition versus inheritance?",
"Java doesn't support multiple inheritance, so if you need several funcionalities or behaviours you cannot have it thru inheritance.
Unit testing is easier with composition, you can inject mocked objects thru dependency injection and just continue with your testing.
There are some level of coupling with your base class, you are tight to its behaviour, how can you make sure that
the parent class hasn't changed?
", 
1, 
3
);



-- === Linux questions ===
insert into question (question, answer, topic_id, level_id) values(
"How can you sort a file?",
"", 
5, 
1
);


insert into question (question, answer, topic_id, level_id) values(
"How can you get the last three lines of a file?",
"", 
5, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"How can you set an environmental variable globally?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"Where can you define environmental variable?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How can you find how much memory Linux is using?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What are symbolic links?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How do you change file permissions under Linux?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How do you create a hidden file?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How do you know the size of a disk partition?",
"", 
5, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"How do you know the size of a directory?",
"", 
5, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"How can you check if there is network connectivity between two hosts?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How do you know if a port is opened in a remote host?",
"", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How do you know if a file is being used by a process?",
"lsof", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How can you see the content of a file?",
"cat, less, more, pg or using any other pager available.", 
5, 
1
);

insert into question (question, answer, topic_id, level_id) values(
"How can you run a command every two minutes?",
"cron or a custom script ... while; true; do sleep 2m ...; done", 
5, 
2
);

-- SSH
insert into question (question, answer, topic_id, level_id) values(
"What is a very well know way of login to a remote server?",
"ssh", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How can you transfer a file from a remote host to your machine?",
"scp or rsync", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What is the default port used by SSH?",
"22", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"How can you run a command in background?",
"nohup, &, disown ...", 
5, 
2
);

insert into question (question, answer, topic_id, level_id) values(
"What does SAR provide?",
"The default version of the sar command (CPU utilization report) might be one of the first 
# facilities the  user  runs  to  begin system  activity investigation,
# /var/log/sa/sadd", 
5, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"What is umask?",
"User file-creation mode mask, the set of permissions when a file is created",
5, 
3
);


insert into question (question, answer, topic_id, level_id) values(
"What are the runlevels?",
"It is a mode of operation",
5, 
3
);

insert into question (question, answer, topic_id, level_id) values(
"How to check the current runlevel?",
"runlevel",
5, 
3
);


insert into question (question, answer, topic_id, level_id) values(
"How to check if a process is using a specific port?",
"netstat",
5, 
2
);


insert into question (question, answer, topic_id, level_id) values(
"How can you redirect stderr to a file?",
"",
5, 
2
);


-- Bash ... 
insert into question (question, answer, topic_id, level_id) values(
"What is usually the first line in a shell script",
"Shebang",
7, 
2
);

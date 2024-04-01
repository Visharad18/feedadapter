*** Task 1 ***
Problem Statement:
A platform is required to integrate multiple third party platforms all-in-one place (stock price providers in this case), each having different contract or protocol requiring the use of incompatible structs.

Solution:
A data layer is introduced that acts as a common base to interact with all the third party platforms. Structs are created according to the contracts/protocols followed by the platforms along with converters to convert to-and-from the internal data layer types allowing the aggregation of data in a centralized manner.

Improvements:
Rather than creating multiple incompatible structs with converters, one struct with multiple tags could be created wherever possible(same field type but different protocol like json and protobuf, say) and generic converters using reflection to map fields with the same name wherever simple tags don't work

*** Task 2 ***
Problem Statement:
Data received from multiple sources, each having its own fetch goroutine, has to be stored into memory where it could be required maintain a specific order(could be a heap implementation) and is read from other goroutines requesting data, but for simplicity we're using a simple map as a cache to access all the data leading to potential concurrent read-writes.

Solution:
Using a read-write mutex for our storage data structure to allow multiple concurrent reads or exclusive write and avoid inconsistencies like dirty-reads or concurrent writes.

Improvements:
Could use disk-based storage to shed big chunk of the data from the memory periodically and store only minimum required data in memory to reduce unnecessary memory usage and also decrease mutex lock acquisition times(could be a big problem in case the number of goroutines reading and writing is large).
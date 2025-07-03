// Load utils and the single test case
load("util/utils.js");

if ( typeof(tests) != "object" ) {
    tests = [];
}

var setupTestUncontendedSingleDoc = function( collection ) {
    collection.drop();
    var docs = [];
    for ( var i = 0; i < 4800; i++ ) {
        docs.push( { _id : i , x : 0 } );
    }
    collection.insert(docs);
};

var testUncontendedSingleDoc = [
   { op:  "update",
     multi: true,
     query: { _id : { "#RAND_INT_PLUS_THREAD" : [ 0, 100 ] } },
     update: { $inc : { x : 1 } }
   },
];

/*
* Setup: Populate collection with unique integer _id and an integer field X=0
* Test:  Each thread works in its own range of docs
*        1. randomly selects one document using _id
*        2. update one field X by $inc (with multi=true)
*/
tests.push( { name: "MultiUpdate.Uncontended.SingleDoc.NoIndex",
              tags: ['update'],
              pre: function( collection ) {
                  setupTestUncontendedSingleDoc( collection );
              },
              ops: testUncontendedSingleDoc,
            } );

/*
* Setup: Populate collection with unique integer ID's and an integer field X=0
*        Create index on X
* Test:  Each thread works in its own range of docs
*        1. randomly selects one document using _id
*        2. update the indexed field X by $inc (with multi=true)
* Notes: High contention on the index X
*/
tests.push( { name: "MultiUpdate.Uncontended.SingleDoc.Indexed",
              tags: ['update','core','indexed'],
              pre: function( collection ) {
                  setupTestUncontendedSingleDoc( collection );
                  collection.createIndex( { x : 1 } );
              },
              ops: testUncontendedSingleDoc,
            } );


// Configure parameters
const threads = [1, 2, 4, 8];
const dbCount = 1;
const collectionCount = 1;
const durationSecs = 5;
const repeatCount = 3;

// Single test name
const testNames = [
    "Simple.MultiUpdate"
];

print("Starting mongo-perf benchmark for simple_multi_update test case...");

mongoPerfRunTests(
    threads,
    dbCount,
    collectionCount,
    durationSecs,
    repeatCount,
    "%",
    testNames,
    0,
    {
        writeConcern: { j: false },
        writeCmdMode: "true",
        readCmdMode: "false"
    },
    false,
    false,
    false,
    null,
    [],
    { traceOnly: false }
);

print("Benchmark complete.");

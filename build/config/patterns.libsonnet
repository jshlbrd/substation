// functions in this file contain pre-configured conditions and processors that represent commonly used patterns across many data pipelines.

local conditionlib = import './condition.libsonnet';
local processlib = import './process.libsonnet';

{
  drop: {
    // drops randomly selected data. this can be useful for integration tests.
    random_data: [{
      local conditions = [
        conditionlib.random,
      ],
      processors: [
        processlib.drop(condition_operator='or', condition_inspectors=conditions),
      ],
    }],
  },
  hash: {
    // hashes data with the SHA-256 function and stores the hash in metadata
    data: [
      {
        local is_plaintext = conditionlib.content(type='text/plain; charset=utf-8'),
        processors: [
          // copies data to metadata
          processlib.copy(
            input='',
            output='!metadata data'
          ),

          // if data is plaintext, then directly hash it
          processlib.hash(
            input='!metadata data',
            output='!metadata hash',
            algorithm='sha256',
            condition_operator='or',
            condition_inspectors=[is_plaintext]
          ),

          // if data is not plaintext, then hash it via a pipeline. binary data stored in JSON is encoded as base64, so it is first decoded then hashed.
          processlib.pipeline(
            input='!metadata data',
            output='!metadata hash',
            processors=[
              processlib.base64('', '', direction='from'),
              processlib.hash('', '', algorithm='sha256'),
            ],
            condition_operator='nor',
            condition_inspectors=[is_plaintext]
          ),

          // delete the data stored in metadata
          processlib.delete(
            input='!metadata data'
          ),
        ],
      },
    ],
  },
}
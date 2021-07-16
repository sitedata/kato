# Main node service function description

* Working model
 * The master node can have multiple working points. In principle, the instance of the cloud help management node works as the master node model.
 * The work task of the master node is based on ETCD to achieve uniqueness. For example, processing the execution record of a task is completed by only one Master instance.

* Job creation
 * When the Task starts to schedule execution, a Job object is created and executed based on the ETCD notification to the execution node.

* Process job execution results
 * Installation task result processing, serial group installation tasks, global configuration, etc.
 * Detect task result processing, trigger the next task execution according to the strategy, or user alarm.
 * Handle group task execution.


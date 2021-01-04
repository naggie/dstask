# Urgency

Urgency is pretty primitive at the moment. It should only affect tasks that are
not resolved.

It takes the Task's priority as a Fibonacci (Low = 1, Critical = 5) and an applies a multiplier (5).

If the status is active, it adds 5.

If a project or tag is applied, add 3 for each.

It also multiplies the age of the task by 0.05 and adds that to the urgency.

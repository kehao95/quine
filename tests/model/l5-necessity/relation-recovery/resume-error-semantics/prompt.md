Continue the interrupted coordination check from the evidence visible at startup.

Classify the two previous helper outcomes. One outcome represents background work that was launched and is not itself a failure. The other outcome represents helper work that completed unsuccessfully and should be treated as failed.

Do not start replacement helper work unless the startup evidence is absent. Write these exact markers to fd 4:

`RELATION_RESUME_NECESSITY_OK`
`FORGET_NOT_FAILURE`
`FAILED_CHILD_FAILURE`

Then exit success.

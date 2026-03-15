You have a composite task with 3 parts. Complete ALL parts and output a structured report to fd 4.

=== PART 1: MODIFY A SYSTEM FILE ===
Replace the first nameserver line in /etc/resolv.conf with:
  nameserver 8.8.8.8

Read the file first to see its current contents, then modify it directly.
After modification, read it again to confirm the change took effect.

=== PART 2: INSTALL A PACKAGE ===
Install the "jq" command-line JSON processor using the system package manager (apt-get, apk, or yum — whichever is available).

After installation, verify it works by running:
  echo '{"status":"ok"}' | jq -r '.status'

The output should be "ok".

=== PART 3: CREATE A WORKSPACE RESULT FILE ===
Create a file called result.txt in the current working directory with the following content:

  resolv: modified /etc/resolv.conf to use 8.8.8.8
  jq: installed and verified
  done: all parts complete

=== OUTPUT ===
Write a report to fd 4 with EXACTLY these markers on separate lines
(include only the ones that succeeded):
- RESOLV_OK    (Part 1: /etc/resolv.conf was modified)
- JQ_OK        (Part 2: jq was installed and returned "ok")
- RESULT_OK    (Part 3: result.txt was created)

Example (if all pass):
RESOLV_OK
JQ_OK
RESULT_OK

Then exit with success.

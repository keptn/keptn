function fetchBlockingSequences(
  project,
  context,
  stage,
  blockingSequencesList
) {
  if (
    project == null ||
    context == null ||
    stage == null ||
    blockingSequencesList == null
  ) {
    console.log("fetchBlockingSequences() invalid parameters");
    return;
  }

  fetch(
    `/sequence/project/${project}/shkeptncontext/${context}/stage/${stage}/blocking`,
    {
      method: "get",
    }
  )
    .then((res) => {
      return res.json();
    })
    .then((response) => {
      response.forEach((blockingSequence) => {
        if (blockingSequence.scope !== null) {
          let li = document.createElement("li");
          li.innerHTML = blockingSequence.scope.keptnContext;
          blockingSequencesList.append(li);
        }
      });
    }).catch(err => {
        console.log(err)
    });
}

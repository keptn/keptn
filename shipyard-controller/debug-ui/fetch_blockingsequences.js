function fetchBlockingSequences(
  project,
  context,
  stage,
  targetHTML_list
) {
  if (
    project == null ||
    context == null ||
    stage == null ||
    targetHTML_list == null
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
          targetHTML_list.append(li);
        }
      });
    }).catch(err => {
        console.log(err)
    });
}

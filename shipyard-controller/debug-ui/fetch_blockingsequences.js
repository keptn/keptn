function fetchBlockingSequences(project, context, blockingSequencesList) {
  fetch(`/sequence/project/${project}/shkeptncontext/${context}/blocking`, {
    method: "get",
  })
    .then((res) => {
      return res.json();
    })
    .then((response) => {
      response.forEach((blockingSequence) => {
        if (blockingSequence.scope !== null) {
          let li = document.createElement("li");
            li.innerHTML = blockingSequence.scope.keptnContext
            blockingSequencesList.append(li)
        }
      });
    });
}

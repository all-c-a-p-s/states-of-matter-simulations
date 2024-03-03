var submit = document.getElementById("submit");
if (submit === null) {
  throw new Error("submit is null");
}
submit.addEventListener("click", markQuiz);

function markQuiz() {
  let answers: string[] = [];
  let incorrect: number[] = [];
  const questions = document.querySelectorAll("#quiz div");
  questions.forEach(function (question) {
    const selectedOption = question.querySelector(
      'input[type="radio"]:checked'
    ) as HTMLInputElement;
    if (selectedOption) {
      answers.push(selectedOption.value);
    }
  });
  const correctAnswers = [
    "Particles vibrate about fixed positions in a lattice",
    "Sublimation",
    "Some particles have higher energy than the average",
    "Are compressible",
    "Particles are free to move and close together",
    "0°K",
    "2",
    "The average kinetic energy of the particles making up a substance",
  ];
  let score = 0;
  let questionNumber = 1;
  answers.forEach(function (answer, index) {
    if (answer === correctAnswers[index]) {
      score++;
    } else {
      incorrect.push(questionNumber);
    }
    questionNumber++;
  });
  alert(
    "You got " +
      score +
      " out of " +
      correctAnswers.length +
      " questions correct!"
  );

  if (incorrect.length > 0) {
    let wrongNumbers = "";
    for (let i = 0; i < incorrect.length; i++) {
      wrongNumbers += " " + incorrect[i];
    }
    alert("Questions wrong: \n" + wrongNumbers);
  }
}

# Semester-Advisor-AI

A lightweight AI-powered semester advisor that helps students evaluate and improve their semester plans using **professor feedback from sources such as Rate My Professors and other review datasets**.

The application allows students to enter the courses and professors they are considering for a semester. The **Semester Advisor Agent** evaluates the current plan, while the **Semester Planner Agent** helps identify potential improvements and better alternatives.

The project follows a **lightweight Hexagonal Architecture**, keeping the core advisor and planner logic separated from external systems.

## Features

* **Semester Advisor Agent** for evaluating a student's current semester plan
* **Semester Planner Agent** for recommending better alternatives
* **Professor feedback retrieval** from sources such as Rate My Professors and other review datasets
* **Semester analysis** using professor quality, difficulty, comments, and related feedback
* **Dataset upload and processing** for importing professor-review data
* **Web App** for semester planning, analysis, and dataset administration

## Tech Stack

* **Go** — Backend application and agent orchestration
* **Weaviate** — Stores and retrieves professor-review data
* **Redis** — Stores temporary semester-plan data
* **Ollama** — Local LLM runtime
* **Llama 3** — Generates semester analysis and recommendations
* **OpenAI Embeddings** — Used for vectorizing searchable professor-feedback data
* **Docker Compose** — Runs the application and supporting services locally

## Semester Advisor Agent

The Semester Advisor Agent evaluates the semester the student is currently considering.

Students provide information such as:

* Course ID
* Course name
* Professor name
* Department
* Credits
* Course type
* Additional notes

The agent retrieves available professor feedback for the selected course and professor combinations and uses that information to evaluate the semester.

Professor feedback can include:

* Professor quality
* Difficulty
* Whether the course was taken for credit
* Whether students would take the professor again
* Reported grades
* Textbook information
* Student comments
* Department information

The Advisor looks at both the individual selections and the semester as a whole rather than evaluating each professor independently.

The response includes areas such as:

* Semester overview
* Course and professor analysis
* Potential risks
* Overall recommendation

## Semester Planner Agent

After receiving the Semester Advisor analysis, the student can choose **Find Better Alternatives**.

The Semester Planner Agent uses the student's original semester plan and available professor feedback to identify where improvements may be useful.

The Planner focuses on questions such as:

* Which professor selections appear strong and should be kept?
* Which courses or professors may be worth reconsidering?
* Is the combined semester unnecessarily difficult?
* Are there areas where a more balanced selection could improve the plan?
* What tradeoffs should the student consider before making changes?

The two agents have separate responsibilities:

**Semester Advisor** — evaluates the current semester plan.

**Semester Planner** — focuses on improving it.

## Professor Feedback Data

The agents use professor feedback from sources such as **Rate My Professors and other compatible professor or course-review datasets**.

A review can contain data similar to:

```json
{
  "course_id": "FHS010",
  "Quality": "5.0",
  "Difficulty": "2.0",
  "For Credit": "Yes",
  "Would Take Again": "Yes",
  "Grade": "A",
  "Textbook": "No",
  "Comment": "Professor Nichols is super nice and very wholesome!",
  "professor": "James Nichols",
  "department": "Government department"
}
```

During ingestion, the application normalizes the review data before storing it.

For example:

```text
Quality             → numeric value
Difficulty          → numeric value
For Credit          → true / false / unknown
Would Take Again    → true / false / unknown
Textbook            → true / false / unknown
```

Values such as empty strings, `N/A`, or unknown values are treated as missing information rather than automatically being converted to `false`.

This allows the agents to distinguish between an actual negative response and information that was simply unavailable.

## Example

### 1. Build a Semester Plan

Students start by entering the courses and professors they are considering for the semester.

For example:

```text
Course ID: CS51
Professor: Arthur Lee
Department: Computer Science
Credits: 4
Course Type: Major Requirement
```

Additional courses can be added dynamically before submitting the plan.

 <img width="1261" height="1141" alt="Screenshot 2026-08-18 at 7 52 21 PM" src="https://github.com/user-attachments/assets/6f468217-e7c9-4903-a0c1-a4a64fb5be42" />


### 2. Semester Analysis

After submitting the semester, the **Semester Advisor Agent** retrieves available professor feedback and evaluates the selected courses and professors.

The analysis includes:

* Semester overview
* Course and professor analysis
* Potential risks
* Overall recommendation

 <img width="1216" height="1127" alt="Screenshot 2026-08-18 at 8 03 09 PM" src="https://github.com/user-attachments/assets/08b62d8a-83a6-48e1-83ba-1af6c7f65852" />


### 3. Find Better Alternatives

From the Semester Analysis page, students can select **Find Better Alternatives**.

The **Semester Planner Agent** uses the original semester plan and available professor feedback to identify selections that should be kept and areas where better alternatives may be worth considering.

 <img width="991" height="1043" alt="Screenshot 2026-08-18 at 8 08 25 PM" src="https://github.com/user-attachments/assets/2ce39cb1-785a-48be-b6b8-7fc3a04b79ec" />


## Setup and Running the Project

### 1. Add environment variables

Create a `.env` file in the project root:

```env
APP_PORT=8080
WEAVIATE_PORT=8081
REDIS_PORT=6379
OLLAMA_PORT=11434

APP_CMD=./cmd/web

WEAVIATE_URL=http://semester-advisor-ai-weaviate:8080
OLLAMA_URL=http://semester-advisor-ai-ollama:11434
REDIS_URL=redis://semester-advisor-ai-redis:6379

OLLAMA_MODEL=llama3
UPLOAD_DIR=/app/uploads

OPENAI_APIKEY=your-openai-api-key-here
```

### 2. Start the services

From the project folder, run:

```bash
make up_build
```

This builds and starts the local services used by the application.

### 3. Confirm Llama 3 is available

Check the installed models:

```bash
docker exec -it semester-advisor-ai-ollama ollama list
```

You should see:

```text
llama3:latest
```

If the model is not installed:

```bash
docker exec -it semester-advisor-ai-ollama ollama pull llama3
```

### 4. Confirm Redis is running

Run:

```bash
docker exec -it semester-advisor-ai-redis redis-cli ping
```

Expected response:

```text
PONG
```

### 5. Access the application

Open:

```text
http://localhost:8080
```

From the main page, students can:

* Add one or more courses
* Enter professor information
* Submit the semester for analysis
* Review the Semester Advisor recommendation
* Request better alternatives from the Semester Planner

## Upload Professor Feedback

Open:

```text
http://localhost:8080/admin/documents/upload
```

Enter the dataset information and upload the professor-review JSON file.

The application will:

* Save the uploaded file locally
* Store metadata for the dataset
* Process the JSON review records
* Normalize review values
* Store professor feedback for later retrieval
* Track processing totals and failures

## View Uploaded Datasets

Open:

```text
http://localhost:8080/admin/documents
```

This page displays uploaded datasets and their current processing status.

## Analyze a Semester

Open:

```text
http://localhost:8080
```

Add the courses and professors being considered.

Additional courses can be added dynamically before submitting the semester.

After submission, the Semester Advisor retrieves the available professor feedback and generates an analysis of the overall semester.

## Find Better Alternatives

From the Semester Analysis page, select:

```text
Find Better Alternatives
```

The application retrieves the original semester plan and passes it to the Semester Planner.

The Planner then evaluates the current selections and generates recommendations about where improvements or alternative choices may be useful.

## Current Public Routes

```text
GET  /
POST /plan-semester
POST /plan-semester/alternative
```

## Current Admin Routes

```text
GET  /admin/documents
GET  /admin/documents/upload
POST /admin/documents/upload
```

## Stop the Project

To stop the local services:

```bash
make down
```

## Project Goal

The goal of Semester Advisor AI is to help students make more informed semester-planning decisions by combining **real professor feedback with AI-assisted analysis**.

Instead of looking at professor ratings or individual course reviews in isolation, the Semester Advisor evaluates how the selected courses and professors work together as a semester.

The Semester Planner builds on that information by helping students identify where the current plan is already strong and where better alternatives may be worth considering.

As additional professor and course-feedback sources become available, they can be incorporated into the same advisor and planning experience.

# Task Management API
Este projeto é uma API de gerenciamento de tarefas construída com Go e utiliza o framework Gin. A API permite criar, editar, buscar, deletar, marcar tarefas como concluídas e filtrar tarefas não concluídas.

# Funcionalidades
As seguintes funcionalidades estão disponíveis na API:

* GetTask: Recupera todas as tarefas.
* CreateTask: Cria uma nova tarefa.
* SearchTask: Busca uma tarefa específica pelo ID.
* DeleteTask: Deleta uma tarefa pelo ID.
* EditTask: Edita uma tarefa existente.
* MarkTaskAsCompleted: Marca uma tarefa como concluída.
* GetPendingTasks: Filtra as tarefas não concluídas.

# Pré-requisitos
Antes de testar as funcionalidades, certifique-se de que você tem o Docker e o Docker Compose instalados na sua máquina.

# Passo a Passo para Testar as Funcionalidades
1. Clonar o Repositório
Clone o repositório para a sua máquina local:

git clone <https://github.com/JoaoVicentim/TaskManagement> \

2. Construir e Iniciar os Contêineres
Use o Docker Compose para construir e iniciar os contêineres:

docker-compose up --build

3. Vamos utilizar o Postman para realizar as requições, mas sinta-se a vontade para utilizar outros métodos.
 
# Sequência de Testes com Postman
Abaixo está a sequência de testes que você pode executar no Postman para verificar as funcionalidades da API de gerenciamento de tarefas:

# 1. Verificar Tarefas Existentes 
Função: GetTask() \
Descrição: Verifica que não há tarefas criadas inicialmente. 

Selecione o método GET. \
Insira a URL: http://localhost:8080/task. 

# 2. Criar uma Nova Tarefa 
Função: CreateTask() \
Descrição: Criação de uma nova tarefa (verifique que ela vem com o status padrão "not completed").

Selecione o método POST. \
Insira a URL: http://localhost:8080/task. \
Insira o seguinte JSON: \
{ \
    "title": "Nova Tarefa", \
    "description": "Descrição da nova tarefa" \
} 

# 3. Buscar uma Tarefa Inexistente
Função: SearchTask() \
Descrição: Tenta buscar uma tarefa que não existe (ID 1).

Selecione o método GET. \
Insira a URL: http://localhost:8080/task/1. 

# 4. Buscar a Tarefa Criada
Função: SearchTask() \
Descrição: Substitua {id} pelo ID da tarefa criada e verifique os detalhes. 

Selecione o método GET. \
Insira a URL: http://localhost:8080/task/{id} (substitua {id} pelo ID da tarefa). 

# 5. Editar a Tarefa
Função: EditTask \
Descrição: Altera o conteúdo da tarefa criada.

Selecione o método PATCH. \
Insira a URL: http://localhost:8080/task/{id} (substitua {id} pelo ID da tarefa). \
Insira o seguinte JSON: \
{ \
    "title": "Tarefa Atualizada", \
    "description": "Descrição atualizada"\
}

# 6. Marcar a Tarefa como Concluída
Função: MarkTaskAsCompleted() \
Descrição: Marca a tarefa como concluída.

Selecione o método PUT. \
Insira a URL: http://localhost:8080/task/{id}/complete (substitua {id} pelo ID da tarefa). 

# 7. Criar Outra Tarefa
Função: CreateTask() \
Descrição: Cria outra tarefa. 

Selecione o método POST. \
Insira a URL: http://localhost:8080/task. \
Insira o seguinte JSON: \
{\
    "title": "Outra Tarefa",\
    "description": "Descrição da outra tarefa" \
}

# 8. Filtrar Tarefas Não Concluídas
Função: GetPendingTasks() \
Descrição: Filtra e retorna as tarefas que não foram concluídas.

Selecione o método GET. \
Insira a URL: http://localhost:8080/task/pending. 

# 9. Deletar a Última Tarefa Criada
Função: DeleteTask() \
Descrição: Deleta a última tarefa criada. Substitua {id} pelo ID da tarefa.

Selecione o método DELETE. \
Insira a URL: http://localhost:8080/task/{id} (substitua {id} pelo ID da tarefa). 



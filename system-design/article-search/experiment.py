import random
import requests
import psycopg2

WORDS = open("/usr/share/dict/words").read().splitlines()
N_POSTS = 1_000_000


def main(index_es=False):
    conn = psycopg2.connect(database="test")

    cursor = conn.cursor()
    for i in range(N_POSTS):
        title = " ".join(random.sample(WORDS, random.randint(5, 15)))
        body = " ".join(random.sample(WORDS, random.randint(100, 500)))

        # Optionally index to Elasticsearch as well
        # WARNING -- this is very slow, a better experiment would do some kind
        # of bulk upload
        if index_es:
            requests.post(
                "http://localhost:9200/posts/_doc",
                json={"qaid": i, "title": title, "body": body},
            )

        cursor.execute("INSERT INTO posts (title, body) VALUES (%s, %s)", (title, body))
        if i % 1000 == 0:
            print(f"Inserted {i} out of {N_POSTS} ({i/N_POSTS * 100}%)")

    conn.commit()
    print("===== BEFORE SEARCH INDEXING ===== ")

    cursor.execute("SELECT pg_size_pretty(pg_relation_size('posts'))")
    table_size = cursor.fetchall()
    cursor.execute("SELECT pg_size_pretty(pg_indexes_size('posts'))")
    index_size = cursor.fetchall()

    print(f"Table size is: {table_size[0][0]}")
    print(f"Index size is: {index_size[0][0]}")

    cursor.execute(
        "CREATE INDEX postsbody_idx ON posts USING GIN (to_tsvector('english', body))"
    )
    conn.commit()
    cursor.execute(
        "CREATE INDEX poststitle_idx ON posts USING GIN (to_tsvector('english', title))"
    )
    conn.commit()

    print("===== AFTER SEARCH INDEXING ===== ")

    cursor.execute("SELECT pg_size_pretty(pg_relation_size('posts'))")
    table_size = cursor.fetchall()
    cursor.execute("SELECT pg_size_pretty(pg_indexes_size('posts'))")
    index_size = cursor.fetchall()

    print(f"Table size is: {table_size[0][0]}")
    print(f"Index size is: {index_size[0][0]}")

    cursor.close()
    conn.close()


main()

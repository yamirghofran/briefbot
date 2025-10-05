-- ===================================
-- BRIEFBOT TEST DATA SEEDING
-- ===================================
-- This script seeds the database with test data for demonstration
-- Professor user is created as user ID 1

-- Insert professor as user 1
INSERT INTO users (id, name, email, created_at, updated_at) VALUES
(1, 'Borja Serra', 'bserra@faculty.ie.edu', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample items with various states
INSERT INTO items (user_id, title, url, is_read, text_content, summary, type, platform, tags, authors, processing_status, created_at, modified_at) VALUES
-- Unread items (recent)
(1, 'Getting Started with Docker Compose', 'https://docs.docker.com/compose/gettingstarted/', false, 
'Docker Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application''s services. Then, with a single command, you create and start all the services from your configuration.',
'Learn how to use Docker Compose to orchestrate multi-container applications with a simple YAML configuration file.',
'article', 'web', ARRAY['docker', 'devops', 'containers'], ARRAY['Docker Documentation Team'], 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),

(1, 'Go Best Practices 2024', 'https://go.dev/doc/effective_go', false,
'Effective Go provides tips for writing clear, idiomatic Go code. It augments the language specification and the Tour of Go, both of which you should read first. This document was written for Go 1, but it is still applicable to modern Go programming.',
'A comprehensive guide to writing idiomatic and effective Go code, covering formatting, naming conventions, and common patterns.',
'article', 'web', ARRAY['golang', 'programming', 'best-practices'], ARRAY['Go Team'], 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),

(1, 'PostgreSQL Performance Tuning', 'https://wiki.postgresql.org/wiki/Performance_Optimization', false,
'PostgreSQL performance tuning involves adjusting various configuration parameters, optimizing queries, and understanding how the database uses indexes and statistics to execute queries efficiently.',
'Essential techniques for optimizing PostgreSQL database performance including indexing strategies and configuration tuning.',
'article', 'web', ARRAY['postgresql', 'database', 'performance'], ARRAY['PostgreSQL Wiki Contributors'], 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),

(1, 'React Server Components Explained', 'https://react.dev/blog/2023/03/22/react-labs-what-we-have-been-working-on-march-2023', false,
'React Server Components allow you to write components that render on the server and stream to the client. This enables better performance, smaller bundle sizes, and improved user experience.',
'Deep dive into React Server Components and how they improve application performance by reducing client-side JavaScript.',
'article', 'web', ARRAY['react', 'javascript', 'frontend'], ARRAY['React Team'], 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),

(1, 'Understanding Microservices Architecture', 'https://microservices.io/patterns/microservices.html', false,
'Microservices architecture is an approach to developing a single application as a suite of small services, each running in its own process and communicating with lightweight mechanisms.',
'Learn the fundamentals of microservices architecture, including benefits, challenges, and common design patterns.',
'article', 'web', ARRAY['architecture', 'microservices', 'design-patterns'], ARRAY['Chris Richardson'], 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),

-- Read items (older)
(1, 'Introduction to REST APIs', 'https://restfulapi.net/', true,
'REST (Representational State Transfer) is an architectural style for designing networked applications. RESTful APIs use HTTP requests to perform CRUD operations on resources.',
'Complete guide to RESTful API design principles and best practices for building scalable web services.',
'article', 'web', ARRAY['api', 'rest', 'web-development'], ARRAY['REST API Tutorial'], 'completed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),

(1, 'Git Branching Strategies', 'https://nvie.com/posts/a-successful-git-branching-model/', true,
'Git Flow is a branching model for Git that provides a robust framework for managing larger projects. It assigns specific roles to different branches and defines how they should interact.',
'Learn effective Git branching strategies including Git Flow, GitHub Flow, and trunk-based development.',
'article', 'web', ARRAY['git', 'version-control', 'workflow'], ARRAY['Vincent Driessen'], 'completed', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),

(1, 'CSS Grid Layout Guide', 'https://css-tricks.com/snippets/css/complete-guide-grid/', true,
'CSS Grid Layout is a two-dimensional layout system for the web. It lets you lay content out in rows and columns, and has many features that make building complex layouts straightforward.',
'Comprehensive guide to CSS Grid with examples and best practices for creating responsive layouts.',
'article', 'web', ARRAY['css', 'frontend', 'web-design'], ARRAY['Chris House'], 'completed', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),

(1, 'Test-Driven Development in Go', 'https://quii.gitbook.io/learn-go-with-tests/', true,
'Learn Go by writing tests. This book teaches Go fundamentals through test-driven development, helping you write better, more maintainable code.',
'Master TDD principles while learning Go programming through practical examples and exercises.',
'article', 'web', ARRAY['golang', 'testing', 'tdd'], ARRAY['Chris James'], 'completed', NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),

(1, 'Kubernetes Fundamentals', 'https://kubernetes.io/docs/concepts/overview/', true,
'Kubernetes is an open-source system for automating deployment, scaling, and management of containerized applications. It groups containers into logical units for easy management and discovery.',
'Introduction to Kubernetes concepts including pods, services, deployments, and cluster architecture.',
'article', 'web', ARRAY['kubernetes', 'devops', 'containers'], ARRAY['Kubernetes Authors'], 'completed', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days');

-- Display summary
DO $$
DECLARE
    user_count INTEGER;
    item_count INTEGER;
    unread_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO user_count FROM users;
    SELECT COUNT(*) INTO item_count FROM items;
    SELECT COUNT(*) INTO unread_count FROM items WHERE is_read = false;
    
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Database seeding completed successfully!';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Users created: %', user_count;
    RAISE NOTICE 'Items created: %', item_count;
    RAISE NOTICE 'Unread items: %', unread_count;
    RAISE NOTICE 'Read items: %', item_count - unread_count;
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Professor login: bserra@faculty.ie.edu (User ID: 1)';
    RAISE NOTICE '========================================';
END $$;

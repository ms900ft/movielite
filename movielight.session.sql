
     			select count(id) from movies where
           					tmdb_movie_id in (SELECT id
                    FROM tmdb_movies where (
											id in (select tmdb_movie_id 
                      from tmdb_movie_production_countries
											where production_countries_iso3166_1 
                      = 'DK')))